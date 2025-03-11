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
	id, ok := TcpConnToConnID.Load(conn)
	if !ok {
		logger.Logger.Debug("not found this tcp connection, close directly")
		conn.ForceClose()
		return
	}
	ctx := context.TODO()
	rpc.GetStateClient().ReceiveUplinkMessage(ctx, &pb.StateRequest{
		Endpoint: GetEndpoint(),
		ConnId:   id.(int64),
		Data:     msg,
	})

}

func onConnection(conn *netreactors.TcpConnection) {
	if conn.Connected() {
		//新建连接	
		id := SnowflakeId.Generate()
		info := &ConnInfo{
			ConnID:   id,
			ConnType: ConnTypeTCP,
			TCP:      conn,
		}
		//维持基本连接状态
		TcpConnToConnID.Store(conn, id)
		IDtoConnInfo.Store(id, info)

		logger.Logger.Debug("new connection:", zap.String("addr", conn.PeerAddr().String()))
	} else {
		//关闭连接
		if id, ok := TcpConnToConnID.Load(conn); ok {
			//被动断开：当连接被动断开了(由客户端断开)，执行这里的回调，再调用state提供的清除连接状态 rpc.
			rpc.GetStateClient().ClearConnState(context.TODO(), &pb.StateRequest{
				Endpoint: GetEndpoint(),
				ConnId:   id.(int64),
			})
			//清除基本连接状态
			TcpConnToConnID.Delete(conn)
			IDtoConnInfo.Delete(id)
		}

		logger.Logger.Debug("close connection:", zap.String("addr", conn.PeerAddr().String()))
	}
}

func GetEndpoint() string {
	return configs.GlobalConfig.StateRpcAddr
}
