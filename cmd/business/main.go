package main

import (
	"nbim/configs"
	"nbim/internal/business/api"
	"nbim/pkg/interceptor"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor()))

	//优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	pb.RegisterBusinessExtServer(server, &api.BusinessExtServer{})
	pb.RegisterBusinessIntServer(server, &api.BusinessIntServer{})

	listen, err := net.Listen("tcp", configs.GlobalConfig.BusinessRpcListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("business RPC 服务启动")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
