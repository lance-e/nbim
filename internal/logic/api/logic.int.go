package api

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicIntServer struct {
	pb.UnsafeLogicIntServer
}

// 登陆
func (s *LogicIntServer) ConnSignIn(context.Context, *pb.ConnSignInReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 消息同步
func (s *LogicIntServer) Sync(context.Context, *pb.SyncReq) (*pb.SyncResp, error) {

	return &pb.SyncResp{}, nil
}

// 设备收到消息回执
func (s *LogicIntServer) ReceiveACK(context.Context, *pb.ReceiveACKReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil

}

// 设备离线
func (s *LogicIntServer) Offline(context.Context, *pb.OfflineReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 订阅房间
func (s *LogicIntServer) SubscribeRoom(context.Context, *pb.SubscribeRoomReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 推送
func (s *LogicIntServer) Push(context.Context, *pb.PushReq) (*pb.PushResp, error) {

	return &pb.PushResp{}, nil
}

// 推送消息到房间
func (s *LogicIntServer) PushRoom(context.Context, *pb.PushRoomReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 推送消息到全服
func (s *LogicIntServer) PushAll(context.Context, *pb.PushAllReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 获取设备信息
func (s *LogicIntServer) GetDevice(context.Context, *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {

	return &pb.GetDeviceResp{}, nil
}

// 服务停止
func (s *LogicIntServer) ServerStop(context.Context, *pb.ServerStopReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
