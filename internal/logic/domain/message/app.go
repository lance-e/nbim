package message

import (
	"context"
	"fmt"
	"nbim/pkg/protocol/pb"
)

type app struct{}

var App = new(app)

// 发送好友消息
func (*app) SendMessageToFriend(ctx context.Context, req *pb.SendMessageReq) error {
	return Service.SendToUser(req.UserId, req.DeviceId, int64(req.SessionId), req.Content, req.SendTime)
}

// 发送群组信息
func (*app) SendMessageToGroup(ctx context.Context, req *pb.SendMessageReq) error {
	return Service.SendToGroup(req.UserId, req.DeviceId, req.SessionId, req.Content, req.SendTime)
}

// 消息同步
func (*app) Sync(ctx context.Context, req *pb.SyncReq) (*pb.SyncResp, error) {

	return nil, nil
}

// 设备收到消息回执
func (*app) ReceiveACK(ctx context.Context, req *pb.ReceiveACKReq) error {
	fmt.Printf("message-[%d] ack\n", req.DeviceAck)
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
