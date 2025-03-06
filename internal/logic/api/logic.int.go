package api

import (
	"context"
	"nbim/internal/logic/domain/device"
	"nbim/internal/logic/domain/message"
	"nbim/internal/logic/domain/room"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicIntServer struct {
	pb.UnsafeLogicIntServer
}

// 权限校验
func (s *LogicIntServer) Auth(context.Context, *pb.AuthReq) (*emptypb.Empty, error) {

	return new(emptypb.Empty), nil
}

// 长连接登陆,建立tcp连接时调用
func (s *LogicIntServer) ConnSignIn(ctx context.Context, req *pb.ConnSignInReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, device.App.ConnSignIn(ctx, req)
}

// 消息同步
func (s *LogicIntServer) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {
	return message.App.Sync(ctx, req)
}

// 设备收到消息回执
func (s *LogicIntServer) ReceiveACK(ctx context.Context, req *pb.ReceiveACKReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.ReceiveACK(ctx, req)
}

// 设备离线
func (s *LogicIntServer) Offline(ctx context.Context, req *pb.OfflineReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, device.App.Offline(ctx, req)
}

// 订阅房间
func (s *LogicIntServer) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, room.App.SubscribeRoom(&ctx, req)
}

// 推送
func (s *LogicIntServer) Push(ctx context.Context, req *pb.PushReq) (*pb.PushResp, error) {
	return message.App.Push(ctx, req)
}

// 推送消息到房间
func (s *LogicIntServer) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.PushRoom(ctx, req)
}

// 推送消息到全服
func (s *LogicIntServer) PushAll(ctx context.Context, req *pb.PushAllReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.PushAll(ctx, req)
}

// 获取设备信息
func (s *LogicIntServer) GetDevice(ctx context.Context, req *pb.GetDeviceReq) (*pb.GetDeviceResp, error) {
	return device.App.GetDevice(ctx, req)
}

// 服务停止
func (s *LogicIntServer) ServerStop(ctx context.Context, req *pb.ServerStopReq) (*emptypb.Empty, error) {
	//异步执行服务停止
	go func() {
		err := device.App.ServerStop(ctx, req)
		if err != nil {
			logger.Sugar.Error(err)
		}
	}()
	return &emptypb.Empty{}, nil
}
