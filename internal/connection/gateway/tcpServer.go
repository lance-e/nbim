package gateway

import (
	"context"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
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
	netreactors.Dlog.TurnOnLog()
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
	if value, ok := ConnToInfo.Load(conn); ok && value != nil {
		value.(*ConnInfo).HandleMessage(conn.InBuffer().RetrieveAllString()) //核心逻辑
	}
}

func onConnection(conn *netreactors.TcpConnection) {
	if conn.Connected() {
		//新建连接	
		ConnToInfo.Store(conn, &ConnInfo{
			ConnType: ConnTypeTCP,
			TCP:      conn,
		})
		logger.Logger.Debug("new connection:", zap.String("addr", conn.PeerAddr().String()))
	} else {
		//关闭连接
		value, ok := ConnToInfo.Load(conn)
		if !ok || value == nil {
			return
		}
		info := value.(*ConnInfo)

		logger.Logger.Debug("close connection:", zap.String("addr", conn.PeerAddr().String()), zap.Int64("user_id", info.UserId), zap.Int64("device_id", info.DeviceId))

		//删除tcp连接与对应的对象信息的映射
		ConnToInfo.Delete(conn)
		//删除设备id与对应对象信息的映射
		DeviceToInfo.Delete(info.DeviceId)

		if info.UserId != 0 {
			rpc.GetLogicIntClient().Offline(context.TODO(), &pb.OfflineReq{
				UserId:     info.UserId,
				DeviceId:   info.DeviceId,
				ClientAddr: conn.PeerAddr().String(),
			})
		}
	}
}
