package main

import (
	"fmt"
	"nbim/pkg/logger"
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
	server.SetGoroutineNum(0)

	server.Start()
	ev.Loop()
}

func onMessage(conn *netreactors.TcpConnection, buf *netreactors.Buffer, t time.Time) {
	fmt.Printf("%s\n", buf.RetrieveAllString())
	fmt.Printf("收到消息，我要开始主动断开tcp连接了\n")
	conn.ForceClose()
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
	StartTCPServer("127.0.0.1:8080")
}
