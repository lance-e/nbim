package state

import (
	"nbim/configs"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunMain() {
	//init
	InitCatheState()

	//异步处理rpc请求
	go func() {
		HandleRPC()
	}()

	server := grpc.NewServer() //TODO:UnaryInterceptor

	//优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterStateServer(server, CS.Server)

	listen, err := net.Listen("tcp", configs.GlobalConfig.StateRpcAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("business RPC 服务启动")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
