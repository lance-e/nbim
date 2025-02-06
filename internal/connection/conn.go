package connection

import (
	"container/list"
	"context"
	"nbim/configs"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	netreactors "github.com/lance-e/net-reactors"
	"go.uber.org/zap"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	ConnTypeTCP int8 = 1
	ConnTypeWS  int8 = 2
)

var (
	ConnToInfo   = sync.Map{} //连接与连接对象信息的映射
	DeviceToInfo = sync.Map{} //设备id与连接对象信息的映射
)

type ConnInfo struct {
	ConnType int8
	TCP      *netreactors.TcpConnection
	WSMutex  sync.Mutex
	WS       *websocket.Conn
	UserId   int64
	DeviceId int64
	RoomId   int64
	Element  *list.Element
}

// Write:负责向连接发送数据包
func (info *ConnInfo) Write(data []byte) {
	switch info.ConnType {
	case ConnTypeTCP:
		info.TCP.Send(data)
	case ConnTypeWS:
		//TODO:websocket
	default:
		logger.Logger.Error("unkonwn connection type ", zap.Any("Connection", info))
	}
}

func (info *ConnInfo) HandleMessage(msg []byte) {
	input := new(pb.Input)
	err := proto.Unmarshal(msg, input)
	if err != nil {
		logger.Logger.Error("unmarshal error", zap.Error(err), zap.Int("len", len(msg)))
		return
	}
	logger.Logger.Debug("HandleMessage", zap.Any("input", input))

	//处理未登陆情况
	if input.Type != pb.PackageType_PT_SIGN_IN && info.UserId == 0 {
		//未登陆
		return
	}

	//处理不同类型数据包
	switch input.Type {
	case pb.PackageType_PT_SIGN_IN: //登陆
		info.SignIn(input)
	case pb.PackageType_PT_SYNC: //同步
		info.Sync(input)
	case pb.PackageType_PT_HEARTBEAT: //心跳
		info.HeartBeat(input)
	case pb.PackageType_PT_MESSAGE: //消息转发(目前只有客户端发送ACK这一种情况)
		info.ReceiveMessage(input)
	case pb.PackageType_PT_SUBSCRIBE_ROOM: //订阅房间
		info.SubscribeRoom(input)
	default:
		logger.Logger.Error("unknown package type")
	}

}

// Send:负责将其他格式数据转换成下行数据格式, 响应给客户端 
// 将最初的客户端发送来的上行信息(input),经过相关逻辑处理获得业务响应数据,通过Send 转换成下行信息(output),再响应给客户端
func (info *ConnInfo) Send(t pb.PackageType, requestId int64, message proto.Message, err error) {
	var output = pb.Output{
		Type:      t,
		RequestId: requestId,
	}

	if err != nil {
		s, _ := status.FromError(err)
		output.Code = s.Proto().GetCode()
		output.Message = s.Message()
	}

	if message != nil {
		msgByte, err := proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error("marshal message error", zap.Error(err))
			return
		}
		output.Data = msgByte
	}

	outputByte, err := proto.Marshal(&output)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	//TODO:对于向连接发送消息的成功或失败的处理
	info.Write(outputByte)
}

func (info *ConnInfo) SignIn(input *pb.Input) {
	var signIn pb.SignInInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	//连接层向grpc服务端发送登陆rpc请求
	_, err = rpc.GetLogicIntClient().ConnSignIn(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("request_id", strconv.FormatInt(input.RequestId, 10))), &pb.ConnSignInReq{
		DeviceId:   signIn.DeviceId,
		UserId:     signIn.UserId,
		Token:      signIn.Token,
		ServerAddr: configs.GlobalConfig.ConnectionLocalAddr,
		ClientAddr: info.TCP.PeerAddr().String(),
	})

	//请求成功后，将rpc响应结果发送给客户端
	info.Send(pb.PackageType_PT_SIGN_IN, input.RequestId, nil, err)

	if err != nil {
		return
	}

	//登陆成功
	info.DeviceId = signIn.DeviceId
	info.UserId = signIn.UserId
	DeviceToInfo.Store(info.DeviceId, &info)
}

func (info *ConnInfo) Sync(input *pb.Input) {
	var syncInput pb.SyncInput
	err := proto.Unmarshal(input.Data, &syncInput)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	//连接层向grpc服务端发送消息同步rpc请求
	syncResp, err := rpc.GetLogicIntClient().Sync(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("request_id", strconv.FormatInt(input.RequestId, 10))), &pb.SyncReq{
		DeviceId: info.DeviceId,
		UserId:   info.UserId,
		Seq:      syncInput.Seq,
	})

	//rpc响应成功就返回有数据的下行消息，失败就返回空的下行消息
	var syncOutput = &pb.SyncOutput{}
	if err == nil {
		syncOutput.Message = syncResp.Message
		syncOutput.HasMore = syncResp.HasMore
	}

	//请求成功后，将rpc响应结果发送给客户端
	info.Send(pb.PackageType_PT_SYNC, input.RequestId, syncOutput, err)

}
func (info *ConnInfo) HeartBeat(input *pb.Input) {
	//通过发送空的output下行消息实现心跳检测
	info.Send(pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	logger.Sugar.Infow("HeartBeat", "device_id", info.DeviceId, "user_id", info.UserId)
}

// ReceiveMessage:用于处理客户端向长连接层发送数据包(目前只有ACK一种情况)
func (info *ConnInfo) ReceiveMessage(input *pb.Input) {
	var ack pb.ACK
	err := proto.Unmarshal(input.Data, &ack)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	_, _ = rpc.GetLogicIntClient().ReceiveACK(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("request_id", strconv.FormatInt(input.RequestId, 10))), &pb.ReceiveACKReq{
		UserId:      info.UserId,
		DeviceId:    info.DeviceId,
		ReceiveTime: ack.ReceiveTime,
	})

	//不需要再响应客户端的ACK

}
func (info *ConnInfo) SubscribeRoom(input *pb.Input) {
	var subscribeRoomInput pb.SubscribeRoomInput
	err := proto.Unmarshal(input.Data, &subscribeRoomInput)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	//订阅房间
	SubscribeRoom(info, subscribeRoomInput.RoomId)
	//回执客户端订阅房间的响应(空包)
	info.Send(pb.PackageType_PT_SUBSCRIBE_ROOM, input.RequestId, nil, nil)

	_, err = rpc.GetLogicIntClient().SubscribeRoom(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("request_id", strconv.FormatInt(input.RequestId, 10))), &pb.SubscribeRoomReq{
		UserId:     info.UserId,
		DeviceId:   info.DeviceId,
		RoomId:     subscribeRoomInput.RoomId,
		Seq:        subscribeRoomInput.Seq,
		ServerAddr: configs.GlobalConfig.ConnectionLocalAddr,
	})

	if err != nil {
		logger.Logger.Error("SubscribeRoom error", zap.Error(err))
	}
}
