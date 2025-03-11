package state

import (
	"context"
	"fmt"
	"nbim/pkg/db"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// var CmdChannel chan *cmdContext

const (
	CmdReceiveUplinkMessage = 1
	CmdClearConnState       = 2
)

type cmdContext struct {
	Ctx      *context.Context
	CmdType  uint8
	Endpoint string
	ConnId   int64
	Data     []byte
}

type StateServer struct {
	cmdchannel chan *cmdContext
	pb.UnsafeStateServer
}

// state server：传递上行消息，由长连接网关服务器调用
func (s *StateServer) ReceiveUplinkMessage(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	newctx := context.TODO()
	s.cmdchannel <- &cmdContext{
		Ctx:      &newctx,
		CmdType:  CmdReceiveUplinkMessage,
		Endpoint: req.GetEndpoint(),
		ConnId:   req.GetConnId(),
		Data:     req.GetData(),
	}
	return new(emptypb.Empty), nil
}

// state server：清除连接状态，由长连接网关服务器调用
func (s *StateServer) ClearConnState(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	newctx := context.TODO()
	s.cmdchannel <- &cmdContext{
		Ctx:      &newctx,
		CmdType:  CmdClearConnState,
		Endpoint: req.GetEndpoint(),
		ConnId:   req.GetConnId(),
	}
	return new(emptypb.Empty), nil
}

func (s *StateServer) DelieverDownlinkMessage(ctx context.Context, req *pb.DelieverDownlinkMessageReq) (*emptypb.Empty, error) {
	//获取到指定设备的连接id
	connId, err := db.RedisCli.Get(fmt.Sprint(db.DeviceIdToConnId, req.DeviceId)).Int64()
	if err != nil {
		return &emptypb.Empty{}, err
	}
	//发送下行消息
	for _, down := range req.Downs {
		sendDownlinkMessage(context.TODO(), connId, down)
	}

	return &emptypb.Empty{}, nil
}
