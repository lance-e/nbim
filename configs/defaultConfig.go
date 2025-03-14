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
		/* if err == err.(viper.ConfigFileNotFoundError) { */
		/* panic("without config.yaml") */
		/* } */
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
		GatewayNodeId:        viper.GetInt("gateway_node_id"),
		StateCatheSlotNum:    viper.GetInt("state_cathe_slot_num"),

		ConnectionLocalAddr:     viper.GetString("connection_local_addr"),
		ConnectionTCPListenAddr: viper.GetString("connection_tcp_addr"),
		ConnectionWSListenAddr:  viper.GetString("connection_ws_addr"),
		ConnectionIpconfigAddr:  viper.GetString("connection_ipconfig_addr"),
		ConnectionLogicAddr:     viper.GetString("connection_logic_addr"),
		GatewayRpcAddr:          viper.GetString("gateway_rpc_addr"),
		StateRpcAddr:            viper.GetString("state_rpc_addr"),
		LogicRpcAddr:            viper.GetString("logic_rpc_addr"),
		/* LogicRpcIntAddr:         viper.GetString("logic_rpc_int_addr"), */
		/* LogicRpcExtAddr:         viper.GetString("logic_rpc_ext_addr"), */

		NewLogicIntClient: func() pb.LogicIntClient {
			clientConn, err := grpc.NewClient(viper.GetString("logic_rpc_addr"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(clientConn)
		},

		NewLogicExtClient: func() pb.LogicExtClient {
			clientConn, err := grpc.NewClient(viper.GetString("logic_rpc_addr"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicExtClient(clientConn)
		},

		NewGatewayClient: func() pb.GatewayClient {
			clientConn, err := grpc.NewClient(viper.GetString("gateway_rpc_addr"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewGatewayClient(clientConn)
		},
		NewStateClient: func() pb.StateClient {
			clientConn, err := grpc.NewClient(viper.GetString("state_rpc_addr"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor))
			if err != nil {
				panic(err)
			}
			return pb.NewStateClient(clientConn)
		},
	}
}
