package message

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 发送好友消息
func (*app) SendMessageToFriend(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {

}

// 发送群组信息
func (*app) SendMessageToGroup(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {

}

// 消息同步
func (*app) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {

}

// 设备收到消息回执
func (*app) ReceiveACK(ctx context.Context, req *pb.ReceiveACKReq) (*emptypb.Empty, error) {

}

// 推送
func (*app) Push(ctx context.Context, req *pb.PushReq) (*pb.PushResp, error) {

}

// 推送信息到房间
func (*app) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*emptypb.Empty, error) {

}

// 推送消息到全服
func (*app) PushAll(ctx context.Context, req *pb.PushAllReq) (*emptypb.Empty, error) {

}
