package state

import (
	"context"
	"fmt"
	"nbim/pkg/db"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"time"

	"google.golang.org/protobuf/proto"
)

// 异步处理rpc请求
func HandleRPC() {
	for cmd := range CS.Server.cmdchannel {
		switch cmd.CmdType {
		case CmdReceiveUplinkMessage:
			//
			pbdata := &pb.Data{}
			if err := proto.Unmarshal(cmd.Data, pbdata); err != nil {
				fmt.Printf("protobuf can't unmarshal\n")
			}
			CmdMessageHandler(cmd, pbdata)
		case CmdClearConnState:
			fmt.Printf("客户端登出-%d,清理状态...\n",cmd.ConnId)                   
			fmt.Printf("启动重连定时器，超时将清除所有状态\n")
			CS.connLogout(cmd.Ctx, cmd.ConnId)
		default:
			sendMsg(cmd.ConnId, pb.CMD_Ack, []byte("commmand not defined"))
		}
	}
}

// 处理通过上行链路的收到的消息，解析信令
func CmdMessageHandler(ctx *cmdContext, cmdMessage *pb.Data) {
	switch cmdMessage.Cmd {
	case pb.CMD_Login:
		login(ctx, cmdMessage)
	case pb.CMD_Uplink:
		uplink(ctx, cmdMessage)
	case pb.CMD_Heartbeat:
		heartbeat(ctx, cmdMessage)
	case pb.CMD_Reconn:
		reconn(ctx, cmdMessage)
	case pb.CMD_Ack:
		ack(ctx, cmdMessage)
	}
}

// 登陆
func login(ctx *cmdContext, data *pb.Data) {
	login := pb.LoginMsg{}
	err := proto.Unmarshal(data.GetPayload(), &login)
	if err != nil {
		fmt.Println("login unmarshal failed")
	}
	//状态初始化
	err = CS.connLogin(ctx.Ctx, login.DeviceId, ctx.ConnId)
	if err != nil {
		fmt.Printf("login failed ,error-[%s]\n", err.Error())
		return
	}
	//调用逻辑层设备上线
	_, err = rpc.GetLogicIntClient().ConnSignIn(*ctx.Ctx, &pb.ConnSignInReq{
		DeviceId: login.DeviceId,
		UserId:   login.UserId,
		//TODO:more information(token......)
	})
	if err != nil {
		fmt.Printf("logic server connsignin failed\n")
		return
	}

	//发送登陆ack
	ack := pb.AckMsg{
		Code:     0,
		Message:  "login OK",
		ToType:   pb.CMD_Login,
		ConnId:   ctx.ConnId,
		ClientId: 0,
	}
	payload, err := proto.Marshal(&ack)
	if err != nil {
		fmt.Printf("unmarshal AckMsg failed\n")
		return
	}
	sendMsg(ctx.ConnId, pb.CMD_Ack, payload)
	//TODO:测试
	// sendDownlinkMessage(ctx.Ctx, ctx.ConnId, 0, 0, []byte("welcome!\n"))
	fmt.Printf("login!!!:deviceId-%d,connID-%d\n", login.DeviceId, ctx.ConnId)
}

// 上行消息
func uplink(ctx *cmdContext, data *pb.Data) {
	up := pb.UplinkMsg{}
	err := proto.Unmarshal(data.GetPayload(), &up)
	if err != nil {
		fmt.Println("uplinkMessage unmarshal failed")
	}
	fmt.Printf("%v\n", up)
	fmt.Printf("sessionId - [%d]\n", up.SessionId)
	//比较更新clientID:确保用户维度消息幂等有序
	if CS.compareAndIncrementClientID(ctx.Ctx, ctx.ConnId, up.ClientId, up.SessionId) {
		//先调用业务层的rpc接口，只有rpc返回成功了，才能更新最大消息id和响应客户端成功
		//TODO:调用业务层,在存储后，下发消息
		sendMsgReq := &pb.SendMessageReq{
			UserId:    up.UserId,
			DeviceId:  up.DeviceId,
			SessionId: up.SessionId, //&运算去掉最高位
			Content:   up.UplinkBody,
			SendTime:  time.Now().UnixMicro(),
		}
		//解析session
		if up.SessionId>>63 == 1 {
			//群聊
			_, err = rpc.GetLogicIntClient().SendMessageToGroup(*ctx.Ctx, sendMsgReq)
		} else {
			//单聊
			_, err = rpc.GetLogicIntClient().SendMessageToFriend(*ctx.Ctx, sendMsgReq)
		}
		//等待rpc成功返回才继续发送ack，
		//失败代表分配seqID失败获取落库失败，
		//防止丢消息，需要客户端重传
		if err != nil {
			fmt.Println(err)
			fmt.Printf("message storage database failed\n")
			ack := pb.AckMsg{
				Code:      0,
				Message:   "nak , storage database failed",
				ToType:    pb.CMD_Uplink,
				ConnId:    ctx.ConnId,
				ClientId:  up.ClientId,
				SessionId: up.SessionId,
			}
			payload, err := proto.Marshal(&ack)
			if err != nil {
				fmt.Printf("AckMsg marshal failed\n")
				return
			}
			sendMsg(ctx.ConnId, pb.CMD_Ack, payload)
			return
		}
		//发送上行消息ack
		ack := pb.AckMsg{
			Code:      0,
			Message:   "OK",
			ToType:    pb.CMD_Uplink,
			ConnId:    ctx.ConnId,
			ClientId:  up.ClientId,
			SessionId: up.SessionId,
		}
		payload, err := proto.Marshal(&ack)
		if err != nil {
			fmt.Printf("AckMsg marshal failed\n")
			return
		}
		sendMsg(ctx.ConnId, pb.CMD_Ack, payload)

	} else {
		fmt.Printf("clientID not match\n")
		slot := GetSlot(ctx.ConnId)
		key := fmt.Sprintf(db.MaxClientIDKey, slot, ctx.ConnId, up.SessionId)
		value, _:= db.RedisCli.Get(key).Int() 
		nck := pb.AckMsg{
			Code:      0,
			Message:   fmt.Sprintf("nak ,lastest is  %d\n",value),
			ToType:    pb.CMD_Uplink,
			ConnId:    ctx.ConnId,
			ClientId:  up.ClientId,
			SessionId: up.SessionId,
		}
		payload, err := proto.Marshal(&nck)
		if err != nil {
			fmt.Printf("AckMsg marshal failed\n")
			return
		}
		sendMsg(ctx.ConnId, pb.CMD_Ack, payload)
	}
	fmt.Print("已处理上行消息\n")
}

// 心跳
func heartbeat(ctx *cmdContext, data *pb.Data) {
	hb := pb.HeartbeatMsg{}
	err := proto.Unmarshal(data.GetPayload(), &hb)
	if err != nil {
		fmt.Println("heartbeat unmarshal failed")
		return
	}
	CS.reSetHeartbeatTimer(ctx.ConnId)
	//发送心跳ack ,这里不发送，减少数据通信量
	fmt.Printf("已处理心跳\n")
}

// 重连
func reconn(ctx *cmdContext, data *pb.Data) {
	fmt.Printf("reconn new connID:%d\n", ctx.ConnId)
	rc := pb.ReconnMsg{}
	err := proto.Unmarshal(data.GetPayload(), &rc)
	if err != nil {
		fmt.Println("reconnMsg unmarshal failed")
		return
	}
	fmt.Printf("reconn old ConnID:%d\n", rc.ConnId)
	var code int64
	msg := "reconnect OK"
	err = CS.connReconn(ctx.Ctx, rc.GetConnId(), ctx.ConnId)
	if err != nil {
		fmt.Println("reconnect failed")
		code = 1
		msg = "reconnect failed"
	}
	//发送重连成功ack
	ack := pb.AckMsg{
		Code:     code,
		Message:  msg,
		ToType:   pb.CMD_Reconn,
		ConnId:   ctx.ConnId,
		ClientId: 0,
	}
	payload, err := proto.Marshal(&ack)
	if err != nil {
		fmt.Printf("AckMsg marshal failed\n")
		return
	}
	sendMsg(ctx.ConnId, pb.CMD_Ack, payload)
	fmt.Printf("已处理重连\n")
}

// 处理下行消息的ack回应
func ack(ctx *cmdContext, data *pb.Data) {
	am := pb.AckMsg{}
	err := proto.Unmarshal(data.GetPayload(), &am)
	if err != nil {
		fmt.Println("AckMsg unmarshal failed")
		return
	}
	//业务的ack
	_, err = rpc.GetLogicIntClient().ReceiveACK(*ctx.Ctx, &pb.ReceiveACKReq{
		UserId:      am.UserId,
		DeviceId:    am.DeviceId,
		DeviceAck:   am.MessageId,
		ReceiveTime: time.Now().UnixMilli(),
	})
	if err != nil {
		logger.Logger.Debug("ReceiveACK rpc failed")
		return
	}
	//状态的ack
	CS.connAck(ctx.Ctx, ctx.ConnId, am.SessionId, am.MessageId)

	fmt.Printf("已处理下行消息ack\n")
}

// 发送下行实体消息
// TODO:(？单聊情况下：这里的sessionID是不是应该传发送方的用户id)
func sendDownlinkMessage(ctx context.Context, connID int64, down *pb.DownlinkMsg) {
	//推送下行实体消息
	payload, err := proto.Marshal(down)
	if err != nil {
		fmt.Printf("DownlinkMsg marshal failed\n")
		return
	}
	sendMsg(connID, pb.CMD_Downlink, payload)
	//------------------------------------------
	//更新状态保存的最后发送的消息
	err = CS.AppendLastMsg(ctx, connID, down)
	if err != nil {
		//TODO:重启state后，内存中的状态信息会消失,需要兜底
		fmt.Println(err)
	}
}

func reSendDownlinkMessage(connID int64) {
	down, err := CS.GetLastMsg(context.Background(), connID)
	if err != nil {
		fmt.Printf("reSendDownlinkMessage error:%s\n", err.Error())
	}
	if down == nil {
		return
	}
	downData, err := proto.Marshal(down)
	if err != nil {
		fmt.Printf("reSendDownlinkMessage unmarshal failed\n")
		return
	}
	sendMsg(connID, pb.CMD_Downlink, downData)
	//重置消息定时器
	if state, ok := CS.ConnIdToConnState.Load(connID); ok {
		state.(*connState).SetMessageTimer()
	}
}

// 发送消息
func sendMsg(connID int64, cmd pb.CMD, payload []byte) {
	data := pb.Data{
		Cmd:     cmd,
		Payload: payload,
	}
	d, err := proto.Marshal(&data)
	if err != nil {
		fmt.Printf("state handler sendMsg: protobuf Marshal failed\n")
	}
	fmt.Printf("下行:data-[%s]\n", string(payload))
	_, err = rpc.GetGatewayClient().SendDownlinkMessage(context.TODO(), &pb.GatewayRequest{
		ConnId: connID,
		Data:   d,
	})
	if err != nil {
		fmt.Printf("GetGatewayClient().SendDownlinkMessage failed\n")
	}
}
