package ipconfig

import (
	"context"
	"math"
	"nbim/pkg/discovery"
	"sort"
	"sync"
)

//etcd中关于ipconfig存的是 ipconfig/node号 : endpoint序列化数据

// dispatcher 的table 存的是 ip:port , endpoint
type dispatcher struct {
	EndPointTable map[string]*discovery.Endpoint
	sync.Mutex
}

var Dispatcher dispatcher

func init() {
	Dispatcher.EndPointTable = make(map[string]*discovery.Endpoint)
}

func Dispatch(ctx context.Context) []*discovery.Endpoint {
	//step1: 获取候选endpoints
	eps := make([]*discovery.Endpoint, 0, len(Dispatcher.EndPointTable))
	for _, ep := range Dispatcher.EndPointTable {
		eps = append(eps, ep)
	}
	//step2: 依次计算得分
	//对于网关机来说，存在不同时期进入的机器，机器配置也是不同的，使用单纯的负载均衡会导致误差较大。我们采用机器的剩余资源作为衡量指标
	//这里评判的是机器剩余可连接数(connect_num) 和 每秒剩余可发送字节数(message_byte) ，message_byte 用于计算活跃分数，connect_num 用于计算静态分数
	activeScore := make([]float64, len(eps))
	staticScore := make([]float64, len(eps))
	for i, ep := range eps {
		GB := ep.MetaData["message_byte"].(float64) / (1 << 30)
		activeScore[i] = math.Trunc(GB*1e2+0.5) * 1e-2
		staticScore[i] = ep.MetaData["connect_num"].(float64)
	}

	//step3: 全局排序，动静结合的排序策略
	//机器剩余资源越多，负载越小
	sort.Slice(eps, func(i, j int) bool {
		if activeScore[i] > activeScore[j] {
			return true
		} else if activeScore[i] == activeScore[j] && staticScore[i] > staticScore[j] {
			return true
		} else {
			return false
		}
	})
	return eps
}
