package logichttp

import (
	"nbim/configs"
)

func RunMain() {

	//http服务
	engine := InitRouter()

	engine.Run(configs.GlobalConfig.ConnectionLogicAddr)
}
