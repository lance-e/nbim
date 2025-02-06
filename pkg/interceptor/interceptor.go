package interceptor

import (
	"context"
	"nbim/pkg/logger"

	"google.golang.org/grpc"
)

// NewInterceptor 生成服务端一元拦截器
func NewInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger.Logger.Debug("TODO: implement NewInterceptor\n")
		//https://developer.aliyun.com/article/1152363
		return
	}
}
