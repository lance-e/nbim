package friend

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 添加好友
func (*app) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*emptypb.Empty, error) {

}

// 同意添加好友
func (*app) AgreeFriend(ctx context.Context, req *pb.AgreeFriendReq) (*emptypb.Empty, error) {

}

// 设置好友信息
func (*app) SetFriend(ctx context.Context, req *pb.SetFriendReq) (*pb.SetFriendResp, error) {

}

// 获取所有好友
func (*app) GetAllFriends(ctx context.Context, req *emptypb.Empty) (*pb.GetAllFriendResp, error) {

}
