package configs

import (
	"context"
	"nbim/pkg/gerror"
	"nbim/pkg/protocol/pb"
	"os"
	"time"

	"google.golang.org/grpc"
)

var Builders = map[string]Buildr{
	"default": &DefaultConfig{},
}

var GlobalConfig Configuration

type Buildr interface {
	Build() Configuration
}

type Configuration struct {
	Mysql                string
	RedisHost            string
	RedisPassword        string
	EtcdEndpoints        []string
	EtcdTimeout          time.Duration
	PushRoomSubscribeNum int
	PushAllSubscribeNum  int

	ConnectionLocalAddr     string
	ConnectionRPCListenAddr string
	ConnectionTCPListenAddr string
	ConnectionWSListenAddr  string
	LogicRPCListenAddr      string
	BusinessRpcListenAddr   string
	FileHTTPListenAddr      string

	NewConnectionIntClient func() pb.ConnectionIntClient
	NewLogicIntClient      func() pb.LogicIntClient
	NewBusinessIntClient   func() pb.BusinessIntClient
}

func init() {
	env := os.Getenv("NBIM")
	builder, ok := Builders[env]
	if !ok {
		builder = new(DefaultConfig)
	}
	GlobalConfig = builder.Build()
}

// 客户端一元拦截器
func interceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return gerror.WrapRPCError(err)
}
