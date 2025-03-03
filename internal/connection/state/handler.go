package state

import (
	"context"
	"fmt"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"

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
			panic("commmand not defined")
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
	default:
		panic("unsupport cmd ")
	}
}

// 登陆
func login(ctx *cmdContext, data *pb.Data) {
	login := pb.LoginMsg{}
	err := proto.Unmarshal(data.GetPayload(), &login)
	if err != nil {
		fmt.Println("login unmarshal failed")
	}
	err = CS.connLogin(ctx.Ctx, login.DeviceId, ctx.ConnId)
	if err != nil {
		fmt.Printf("login failed ,error-[%s]\n", err.Error())
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
	//先调用业务层的rpc接口，只有rpc返回成功了，才能更新最大消息id和响应客户端成功
	if CS.compareAndIncrementClientID(ctx.Ctx, ctx.ConnId, up.ClientId, up.SessionId) {
		//发送上行消息ack
		ack := pb.AckMsg{
			Code:     0,
			Message:  "OK",
			ToType:   pb.CMD_Uplink,
			ConnId:   ctx.ConnId,
			ClientId: up.ClientId,
		}
		payload, err := proto.Marshal(&ack)
		if err != nil {
			fmt.Printf("AckMsg marshal failed\n")
			return
		}
		sendMsg(ctx.ConnId, pb.CMD_Ack, payload)

		//TODO:调用业务层,在存储后，下发消息
		CS.ConnIdToConnState.Range(func(key, value any) bool {
			// 给每一个连接发
			sendDownlinkMessage(ctx.Ctx, key.(int64), CS.MessageId, up.SessionId, up.UplinkBody)
			return true
		})

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
	CS.connAck(ctx.Ctx, ctx.ConnId, am.SessionId, am.MessageId)
	fmt.Printf("已处理下行消息ack\n")
}

// 发送下行实体消息
// TODO:
func sendDownlinkMessage(ctx *context.Context, connID, msgId, sessionId int64, body []byte) {
	//推送下行实体消息
	down := pb.DownlinkMsg{
		MessageId:    msgId,
		SessionId:    sessionId,
		DownlinkBody: body,
	}
	payload, err := proto.Marshal(&down)
	if err != nil {
		fmt.Printf("DownlinkMsg marshal failed\n")
		return
	}
	sendMsg(connID, pb.CMD_Downlink, payload)
	//------------------------------------------
	//更新状态保存的最后发送的消息
	err = CS.AppendLastMsg(ctx, connID, &down)
	if err != nil {
		panic(err)
	}
}

func reSendDownlinkMessage(connID int64) {
	down, err := CS.GetLastMsg(context.Background(), connID)
	if err != nil {
		fmt.Print("reSendDownlinkMessage error:%s\n", err.Error())
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
	rpc.GetGatewayClient().SendDownlinkMessage(context.TODO(), &pb.GatewayRequest{
		ConnId: connID,
		Data:   d,
	})
}
