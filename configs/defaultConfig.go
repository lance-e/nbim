package configs

import (
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 默认配置
type DefaultConfig struct{}

func (d *DefaultConfig) Build() Configuration {
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	return Configuration{
		Mysql:                "lance:123456@tcp(127.0.0.1:3306)/nbim?charset=utf8&parseTime=true&loc=Local",
		RedisHost:            "127.0.0.1:6379",
		RedisPassword:        "123456",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectionLocalAddr:     "127.0.0.1:8000",
		ConnectionRPCListenAddr: "127.0.0.1:8000",
		ConnectionTCPListenAddr: "127.0.0.1:8001",
		ConnectionWSListenAddr:  "127.0.0.1:8002",
		LogicRPCListenAddr:      "127.0.0.1:8010",
		BusinessRpcListenAddr:   "127.0.0.1:8020",
		FileHTTPListenAddr:      "127.0.0.1:8030",

		NewConnectionIntClient: func() pb.ConnectionIntClient {
			clientConn, err := grpc.NewClient("addrs:///127.0.0.1:8000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewConnectionIntClient(clientConn)
		},
		NewLogicIntClient: func() pb.LogicIntClient {
			clientConn, err := grpc.NewClient("addrs:///127.0.0.1:8010", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(clientConn)
		},
		NewBusinessIntClient: func() pb.BusinessIntClient {
			clientConn, err := grpc.NewClient("addrs:///127.0.0.1:8020", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewBusinessIntClient(clientConn)
		},
	}
}
