package discovery

import (
	"context"
	"fmt"
	"nbim/configs"
	"nbim/pkg/logger"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Register struct {
	cli           *clientv3.Client
	leaseId       int64
	key           string
	value         string
	ctx           *context.Context
	keepaliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func NewRegister(ctx *context.Context, key string, endpoint *Endpoint, lease int64) (*Register, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   configs.GlobalConfig.EtcdEndpoints,
		DialTimeout: configs.GlobalConfig.EtcdTimeout,
	})
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	reg := &Register{
		cli:   cli,
		key:   key,
		value: Marshal(endpoint),
		ctx:   ctx,
	}

	//申请租约设置时间 KeepAlive
	if err = reg.PutKeyWithLease(lease); err != nil {
		return nil, err
	}
	return reg, nil
}

// 设置租约
func (r *Register) PutKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := r.cli.Grant(*r.ctx, lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = r.cli.Put(*r.ctx, r.key, r.value, clientv3.WithLease(clientv3.LeaseID(resp.ID)))
	if err != nil {
		return err
	}
	//设置续租，定期发送需求请求
	respChan, err := r.cli.KeepAlive(*r.ctx, clientv3.LeaseID(resp.ID))
	if err != nil {
		return err
	}

	r.leaseId = int64(resp.ID)
	r.keepaliveChan = respChan

	return nil
}

func (r *Register) UpdateValue(val *Endpoint) error {
	value := Marshal(val)
	_, err := r.cli.Put(*r.ctx, r.key, value, clientv3.WithLease(clientv3.LeaseID(r.leaseId)))
	if err != nil {
		return err
	}
	r.value = value

	logger.Logger.Info(fmt.Sprintf("Register UpdateValue leaseID=%d key=%s value=%s , successful!\n", r.leaseId, r.key, r.value))
	return nil
}

func (r *Register) ListenLeaseRespChan() {
	for leaseKeepResp := range r.keepaliveChan {
		logger.Logger.Info(fmt.Sprintf("lease success , leaseID=%d key=%s value=%s , resp:%v!\n", r.leaseId, r.key, r.value, leaseKeepResp))
	}
	logger.Logger.Info(fmt.Sprintf("lease failed!!! leaseID=%d key=%s value=%s \n", r.leaseId, r.key, r.value))
}

func (r *Register) Close() error {
	//撤销租约
	_, err := r.cli.Revoke(*r.ctx, clientv3.LeaseID(r.leaseId))
	if err != nil {
		return err
	}
	logger.Logger.Info(fmt.Sprintf("lease close !!! leaseID:%d , key:%s , value:%s  success\n", r.leaseId, r.key, r.value))
	//关闭
	return r.cli.Close()
}
