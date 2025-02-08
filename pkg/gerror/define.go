package gerror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown         = newError(int(codes.Unknown), "未知错误")
	ErrUnauthorized    = newError(10000, "请重新登陆")
	ErrBadRequest      = newError(10001, "请求参数错误")
	ErrBadCode         = newError(10002, "验证码错误")
	ErrNotInGroup      = newError(10002, "用户非群组成员")
	ErrGroupNotExist   = newError(10002, "群组不存在")
	ErrDeviceNotExist  = newError(10002, "设备不存在")
	ErrAlreadyIsFriend = newError(10002, "对方已经是好友")
	ErrUserNotFound    = newError(10002, "用户未找到")
)

func newError(code int, msg string) error {
	return status.New(codes.Code(code), msg).Err()
}
