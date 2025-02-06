package rpc

import (
	"nbim/configs"
	"nbim/pkg/protocol/pb"
)

// 缓存，减少重复函数调用
var (
	connectionIntClient pb.ConnectionIntClient
	logicIntClient      pb.LogicIntClient
	businessIntClient   pb.BusinessIntClient
)

// 获取某一个connection客户端
func GetConnectionIntClient() pb.ConnectionIntClient {
	if connectionIntClient == nil {
		connectionIntClient = configs.GlobalConfig.NewConnectionIntClient()
	}
	return connectionIntClient
}

// 获取某一个logic客户端
func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		logicIntClient = configs.GlobalConfig.NewLogicIntClient()
	}
	return logicIntClient
}

// 获取某一个business客户端
func GetBusinessIntClient() pb.BusinessIntClient {
	if businessIntClient == nil {
		businessIntClient = configs.GlobalConfig.NewBusinessIntClient()
	}
	return businessIntClient
}
