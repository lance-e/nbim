package rpc

import (
	"context"
	"nbim/pkg/gerror"
	"nbim/pkg/logger"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const (
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxToken     = "token"
	CtxRequestId = "request_id"
)

// 从上下文中获取键对于值
func Get(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	value, ok := md[key]
	if !ok {
		return ""
	}

	return value[0]
}

// 获取ctx的请求id
func GetCtxRequestId(ctx context.Context) int64 {
	idstr := Get(ctx, CtxRequestId)
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// 从ctx获取用户信息，直接返回用户id和设备id
func GetCtxUserInfo(ctx context.Context) (int64, int64, error) {
	var userId int64
	var deviceId int64
	var err error

	userIdstr := Get(ctx, CtxUserId)
	userId, err = strconv.ParseInt(userIdstr, 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerror.ErrUnauthorized
	}
	deviceIdstr := Get(ctx, CtxDeviceId)
	deviceId, err = strconv.ParseInt(deviceIdstr, 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerror.ErrUnauthorized
	}
	return userId, deviceId, nil
}

// 获取ctx的token
func GetCtxToken(ctx context.Context) string {
	return Get(ctx, CtxToken)
}
