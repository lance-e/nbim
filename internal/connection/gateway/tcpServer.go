package gateway

import (
	"context"
	"nbim/configs"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"nbim/pkg/tcp"
	"net/netip"
	"time"

	netreactors "github.com/lance-e/net-reactors"
	"go.uber.org/zap"
)

var (
	ev *netreactors.EventLoop
)

// 启动TCP服务端
func StartTCPServer(addr string) {
	ev = netreactors.NewEventLoop()
	address, err := netip.ParseAddrPort(addr)
	if err != nil {
		panic(err)
	}
	server := netreactors.NewTcpServer(ev, &address, "tcpserver")
	server.SetMessageCallback(onMessage)
	server.SetConnectionCallback(onConnection)

	server.Start()
	ev.Loop()
}

func onMessage(conn *netreactors.TcpConnection, buf *netreactors.Buffer, t time.Time) {
	msg, err := tcp.Unpacking(buf)
	if err != nil {
		logger.Logger.Debug("onMessage: message unpacking failed")
		return
	}
	ctx := context.TODO()
	rpc.GetStateClient().ReceiveUplinkMessage(ctx, &pb.StateRequest{
		Endpoint: configs.GlobalConfig.StateRpcAddr,
		ConnId:   0, //TODO
		Message:  msg,
	})

}

func onConnection(conn *netreactors.TcpConnection) {
	if conn.Connected() {
		//新建连接	

		logger.Logger.Debug("new connection:", zap.String("addr", conn.PeerAddr().String()))
	} else {
		//关闭连接

		logger.Logger.Debug("close connection:", zap.String("addr", conn.PeerAddr().String()))
	}
}
