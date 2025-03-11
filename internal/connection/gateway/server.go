package gateway

import (
	"context"
	"nbim/configs"
	"nbim/pkg/id"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var SnowflakeId *id.Snowflake //用于tcp连接的id生成

func RunMain() {
	//init
	CmdChannel = make(chan *cmdContext, 2048)
	var err error
	SnowflakeId, err = id.NewSnowflake(int64(configs.GlobalConfig.GatewayNodeId))
	if err != nil {
		//FIXME
		panic(err)
	}

	//启动tcp长连接服务器
	go func() {
		StartTCPServer(configs.GlobalConfig.ConnectionTCPListenAddr)
	}()

	//启动websocket长连接服务器
	go func() {
		StartWSServer(configs.GlobalConfig.ConnectionWSListenAddr)
	}()

	//异步处理rpc请求
	go func() {
		AsynHandleRPC()
	}()

	//启动mq
	// StartSubscribe()

	server := grpc.NewServer() //TODO:UnaryInterceptor

	//优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		//这里需要特殊处理，长连接就必须要通知rpc 服务端服务停止
		logger.Logger.Info("server stop start", zap.Any("signal", s))
		rpc.GetLogicIntClient().ServerStop(context.TODO(), &pb.ServerStopReq{
			ServerAddr: configs.GlobalConfig.ConnectionLocalAddr,
		})
		logger.Logger.Info("server stop end", zap.Any("signal", s))
		server.GracefulStop()
	}()

	//注册grpc服务端
	pb.RegisterGatewayServer(server, &GatewayServer{cmdchannel: CmdChannel})

	listen, err := net.Listen("tcp", configs.GlobalConfig.GatewayRpcAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("gateway RPC 服务启动")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
