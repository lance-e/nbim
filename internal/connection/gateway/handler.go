package gateway

import (
	"nbim/pkg/logger"
	"nbim/pkg/tcp"

	"go.uber.org/zap/zapcore"
)

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
	buf := tcp.Packing(ctx.Data)
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
