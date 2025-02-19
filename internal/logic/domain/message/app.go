package message

import (
	"context"
	"nbim/pkg/protocol/pb"
)

type app struct{}

var App = new(app)

// 发送好友消息
func (*app) SendMessageToFriend(ctx context.Context, userId int64, deviceId int64, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	return nil, nil
}

// 发送群组信息
func (*app) SendMessageToGroup(ctx context.Context, userId int64, deviceId int64, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {

	return nil, nil
}

// 消息同步
func (*app) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {

	return nil, nil
}

// 设备收到消息回执
func (*app) ReceiveACK(ctx context.Context, req *pb.ReceiveACKReq) error {
	return nil

}

// 推送
func (*app) Push(ctx context.Context, req *pb.PushReq) (*pb.PushResp, error) {
	return nil, nil

}

// 推送信息到房间
func (*app) PushRoom(ctx context.Context, req *pb.PushRoomReq) error {
	return nil

}

// 推送消息到全服
func (*app) PushAll(ctx context.Context, req *pb.PushAllReq) error {
	return nil

}
