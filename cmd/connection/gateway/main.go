package main

import (
	"context"
	"nbim/configs"
	"nbim/internal/connection/gateway"
	"nbim/pkg/interceptor"
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

func main() {

	//启动tcp长连接服务器
	go func() {
		gateway.StartTCPServer(configs.GlobalConfig.ConnectionTCPListenAddr)
	}()

	//启动websocket长连接服务器
	go func() {
		gateway.StartWSServer(configs.GlobalConfig.ConnectionWSListenAddr)
	}()

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor()))

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

	//启动mq
	gateway.StartSubscribe()

	//注册grpc服务端
	pb.RegisterConnectionIntServer(server, &gateway.ConnectionIntServer{})

	listen, err := net.Listen("tcp", configs.GlobalConfig.ConnectionRPCListenAddr)
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("connection RPC 服务启动")
	err = server.Serve(listen)
	if err != nil {
		logger.Logger.Error("serve error", zap.Error(err))
	}
}
