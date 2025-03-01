package gateway

import (
	"context"
	"nbim/pkg/logger"
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
	Data    []byte
}

type GatewayServer struct {
	cmdchannel chan *cmdContext
	pb.UnsafeGatewayServer
}

// gateway server: 发送下行消息 ,由状态服务器调用
func (g *GatewayServer) SendDownlinkMessage(ctx context.Context, req *pb.GatewayRequest) (*emptypb.Empty, error) {
	logger.Logger.Debug("SendDownlinkMessage rpc")
	newctx := context.TODO()
	g.cmdchannel <- &cmdContext{
		Ctx:     &newctx,
		CmdType: CmdSendDownlinkMessage,
		ConnId:  req.ConnId,
		Data:    req.GetData(),
	}
	return &emptypb.Empty{}, nil
}

// gateway server: 断开指定长连接 ,由状态服务器调用
func (g *GatewayServer) CloseConn(ctx context.Context, req *pb.GatewayRequest) (*emptypb.Empty, error) {
	newctx := context.TODO()
	g.cmdchannel <- &cmdContext{
		Ctx:     &newctx,
		CmdType: CmdCloseConn,
		ConnId:  req.ConnId,
		Data:    req.GetData(),
	}
	return &emptypb.Empty{}, nil
}
