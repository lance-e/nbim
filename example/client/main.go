package client

import (
	"nbim/pkg/logger"
	"net/netip"
	"time"

	netreactors "github.com/lance-e/net-reactors"
	"go.uber.org/zap"
)

var (
	ev *netreactors.EventLoop
)

func StartTCPClient(addr string) {
	ev = netreactors.NewEventLoop()
	address, err := netip.ParseAddrPort(addr)
	if err != nil {
		panic(err)
	}
	client := netreactors.NewTcpClient(ev, &address, "TcpClient")
	client.SetMessageCallback(onMessage)
	client.SetConnectionCallback(onConnection)
}
func onMessage(conn *netreactors.TcpConnection, buf *netreactors.Buffer, t time.Time) {
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

func main() {
	StartTCPClient("127.0.0.1:8080")
}
