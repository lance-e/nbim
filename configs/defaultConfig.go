package configs

import (
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 默认配置
type DefaultConfig struct{}

func (d *DefaultConfig) Build() Configuration {
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	viper.SetConfigFile("./configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if err == err.(viper.ConfigFileNotFoundError) {
			panic("without config.yaml")
		}
		panic(err)
	}

	return Configuration{
		Mysql:                viper.GetString("mysql.name") + ":" + viper.GetString("mysql.password") + "@tcp(" + viper.GetString("mysql.ip") + ":" + viper.GetString("mysql.port") + ")/" + viper.GetString("mysql.dbname") + "?charset=utf8&parseTime=true&loc=Local",
		RedisHost:            viper.GetStringSlice("redis.endpoints")[0],
		RedisPassword:        viper.GetString("redis.password"),
		EtcdEndpoints:        viper.GetStringSlice("etcd.endpoints"),
		EtcdTimeout:          viper.GetDuration("etcd.timeout"),
		PushRoomSubscribeNum: viper.GetInt("push_room_subscribe_num"),
		PushAllSubscribeNum:  viper.GetInt("push_all_subscribe_num"),

		ConnectionLocalAddr:     viper.GetString("connection_local_addr"),
		ConnectionRPCListenAddr: viper.GetString("connection_rpc_addr"),
		ConnectionTCPListenAddr: viper.GetString("connection_tcp_addr"),
		ConnectionWSListenAddr:  viper.GetString("connection_ws_addr"),
		LogicRPCListenAddr:      viper.GetString("logic_rpc_listen_addr"),
		BusinessRpcListenAddr:   viper.GetString("business_rpc_listen_addr"),

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
