package ipconfig

import (
	"context"
	"fmt"
	"math/rand/v2"
	"nbim/pkg/discovery"
	"time"
)

func TestEtcd(ctx *context.Context, ip, port, node string) {
	// 模拟服务发现
	go func() {
		endpoint := discovery.Endpoint{
			Ip:   ip,
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":  float64(rand.Int64N(123123412341234124)),
				"message_byte": float64(rand.Int64N(132412341234)),
			},
		}
		reg, err := discovery.NewRegister(ctx, fmt.Sprintf("%s/%s", "ipconfig", node), &endpoint, time.Now().Unix())
		if err != nil {
			panic(err)
		}

		//监听变化
		go reg.ListenLeaseRespChan()

		for {
			endpoint := discovery.Endpoint{
				Ip:   ip,
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":  float64(rand.Int64N(123123412341234124)),
					"message_byte": float64(rand.Int64N(132412341234)),
				},
			}

			reg.UpdateValue(&endpoint)
			time.Sleep(1 * time.Second)
		}
	}()
}
