package state

import (
	"context"
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
