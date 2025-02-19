package ipconfig

import (
	"context"
	"fmt"
	"nbim/pkg/discovery"
	"nbim/pkg/logger"
)

// 先启动etcd服务发现
func DataHandler(ctx *context.Context) {
	dis := discovery.NewDiscovery(ctx)
	defer dis.Close()

	setFunc := func(key, value string) {
		Dispatcher.Lock()
		val, err := discovery.Unmarshal([]byte(value))
		if err != nil {
			logger.Logger.Error(err.Error())
		} else {
			Dispatcher.EndPointTable[fmt.Sprintf("%s:%s", val.Ip, val.Port)] = val
		}
		Dispatcher.Unlock()
	}

	delFunc := func(key, value string) {
		Dispatcher.Lock()
		val, err := discovery.Unmarshal([]byte(value))
		if err != nil {
			logger.Logger.Error(err.Error())
		} else {
			delete(Dispatcher.EndPointTable, fmt.Sprintf("%s:%s", val.Ip, val.Port))
		}
		Dispatcher.Unlock()
	}

	//监听前缀为ipconfig 的kv
	err := dis.WatchService("ipconfig", setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
