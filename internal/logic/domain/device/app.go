package device

import (
	"context"
	"nbim/pkg/protocol/pb"
)

type app struct{}

var App = new(app)

// 注册设备
func (*app) RegisterDevice(ctx context.Context, req *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	deviceId, err := Service.Register(ctx, req)
	return &pb.RegisterDeviceResp{
		DeviceId: deviceId,
	}, err
}

// 登陆
func (*app) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) error {
	return Service.ConnSignIn(ctx, req)
}

// 设备离线
func (*app) Offline(ctx context.Context, req *pb.OfflineReq) error {
	return Service.Offline(ctx, req)
}

// 获取设备信息
func (*app) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {
	device, err := Service.GetDevice(ctx, req)
	return &pb.GetDeviceResp{Device: device}, err
}

func (*app) ListAllOnlineDeviceByUserId(ctx context.Context, userId int64) ([]*pb.Device, error) {
	return Service.ListAllOnlineDeviceByUserId(ctx, userId)
}

// 服务停止
func (*app) ServerStop(ctx context.Context, req *pb.ServerStopReq) error {
	return Service.ServerStop(ctx, req)
}
