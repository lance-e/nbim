package device

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 注册设备
func (*app) RegisterDevice(ctx context.Context, req *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {

}

// 设备离线
func (*app) Offline(ctx context.Context, req *pb.OfflineReq) (*emptypb.Empty, error) {

}

// 登陆
func (*app) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*emptypb.Empty, error) {

}

// 获取设备信息
func (*app) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {

}

// 服务停止
func (*app) ServerStop(ctx context.Context, req *pb.ServerStopReq) (*emptypb.Empty, error) {

}
