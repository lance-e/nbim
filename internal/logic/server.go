package logic

import (
	"nbim/configs"
	"nbim/internal/logic/api"
	"nbim/internal/logic/domain/device"
	"nbim/internal/logic/domain/group"
	"nbim/internal/logic/proxy"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	proxy.DevcieProxy = device.Service
	proxy.GroupProxy = group.Service
}

func RunMain() {
	server := grpc.NewServer() //TODO:UnaryInterceptor

	//优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		logger.Logger.Info("server stop", zap.Any("signal", s))
		server.GracefulStop()
	}()

	//将对外和对内的服务都注册在同一个端口，减少系统复杂度
	pb.RegisterLogicExtServer(server, &api.LogicExtServer{})
	pb.RegisterLogicIntServer(server, &api.LogicIntServer{})

	listen, err := net.Listen("tcp", configs.GlobalConfig.LogicRpcAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("logic RPC 服务启动")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
