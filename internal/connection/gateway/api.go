package gateway

import (
	"context"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/tcp"

	"go.uber.org/zap/zapcore"
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
	logger.Logger.Debug("SendDownlinkMessage rpc")
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

// 异步处理rpc请求
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

// 发送下行消息
func sendDownlinkMessage(ctx *cmdContext) {
	buf := tcp.Packing(ctx.Payload)
	if info, ok := IDtoConnInfo.Load(ctx.ConnId); ok {
		conninfo := info.(*ConnInfo)
		switch conninfo.ConnType {
		case ConnTypeTCP:
			conninfo.TCP.SendFromBuffer(buf)
		case ConnTypeWS:
			//TODO
		default:
			logger.Logger.Debug("unknown connection type", zapcore.Field{})
		}
	}
}

// 关闭连接
func closeConn(ctx *cmdContext) {
	if info, ok := IDtoConnInfo.Load(ctx.ConnId); ok {
		conninfo := info.(*ConnInfo)
		switch conninfo.ConnType {
		case ConnTypeTCP:
			//FIXME:Shutdown()，只关闭写端，保证数据不丢失
			conninfo.TCP.ForceClose() //强制关闭读端写端 
		case ConnTypeWS:
			//TODO
		default:
			logger.Logger.Debug("unknown connection type")
		}
	}
}
