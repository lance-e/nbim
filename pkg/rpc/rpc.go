package rpc

import (
	"nbim/configs"
	"nbim/pkg/protocol/pb"
)

// 单例模式，延迟缓存
// 缓存，减少重复函数调用
var (
	gatewayClient  pb.GatewayClient
	stateClient    pb.StateClient
	logicIntClient pb.LogicIntClient
)

func GetGatewayClient() pb.GatewayClient {
	if gatewayClient == nil {
		gatewayClient = configs.GlobalConfig.NewGatewayClient()
	}
	return gatewayClient
}

func GetStateClient() pb.StateClient {
	if stateClient == nil {
		stateClient = configs.GlobalConfig.NewStateClient()
	}
	return stateClient
}

// 获取某一个logic客户端
func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		logicIntClient = configs.GlobalConfig.NewLogicIntClient()
	}
	return logicIntClient
}
