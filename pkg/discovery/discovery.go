package discovery

import (
	"context"
	"fmt"
	"nbim/configs"
	"nbim/pkg/logger"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discovery struct {
	cli *clientv3.Client
	// lock sync.Mutex
	ctx *context.Context
}

func NewDiscovery(ctx *context.Context) *Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   configs.GlobalConfig.EtcdEndpoints,
		DialTimeout: configs.GlobalConfig.EtcdTimeout,
	})
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	dis := &Discovery{
		cli: cli,
		ctx: ctx,
	}
	return dis
}

// 初始化服务列表和监控
func (d *Discovery) WatchService(prefix string, set func(key, value string), del func(key, value string)) error {
	//根据前缀获取现有的key
	resp, err := d.cli.Get(*d.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	//监听前缀
	watchChan := d.cli.Watch(*d.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision+1))
	logger.Logger.Info(fmt.Sprintf("watching prefix:%s now...\n", prefix))

	for resp := range watchChan {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				set(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}

	}
	return nil
}

func (d *Discovery) Close() error {
	return d.cli.Close()
}
