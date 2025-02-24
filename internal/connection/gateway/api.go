package gateway

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

var CmdChannel chan *cmdContext

const (
	CmdSendDownlinkMessage = 1
	CmdCloseConn           = 2
)

type cmdContext struct {
	Ctx     *context.Context
	CmdType uint8
	ConnId  int64
	Payload []byte
}

type GatewayServer struct {
	cmdchannel chan *cmdContext
	pb.UnsafeGatewayServer
}

// gateway server: 发送下行消息 ,由状态服务器调用
func (g *GatewayServer) SendDownlinkMessage(ctx context.Context, req *pb.GatewayRequest) (*emptypb.Empty, error) {
	newctx := context.TODO()
	CmdChannel <- &cmdContext{
		Ctx:     &newctx,
		CmdType: CmdSendDownlinkMessage,
		ConnId:  req.ConnId,
		Payload: req.GetMessage(),
	}
	return &emptypb.Empty{}, nil
}

// gateway server: 断开指定长连接 ,由状态服务器调用
func (g *GatewayServer) CloseConn(ctx context.Context, req *pb.GatewayRequest) (*emptypb.Empty, error) {
	newctx := context.TODO()
	CmdChannel <- &cmdContext{
		Ctx:     &newctx,
		CmdType: CmdCloseConn,
		ConnId:  req.ConnId,
		Payload: req.GetMessage(),
	}
	return &emptypb.Empty{}, nil
}

func AsynHandleRPC() {
	for cmd := range CmdChannel {
		switch cmd.CmdType {
		case CmdSendDownlinkMessage:
			Product(func() { sendDownlinkMessage(cmd) })
		case CmdCloseConn:
			Product(func() { closeConn(cmd) })
		default:
			panic("commmand not defined")
		}
	}
}

func sendDownlinkMessage(ctx *cmdContext) {

}

func closeConn(ctx *cmdContext) {

}
