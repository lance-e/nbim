package friend

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 添加好友
func (*app) AddFriend(ctx context.Context, userid int64, req *pb.AddFriendReq) error {
	return Service.AddFriend(ctx, userid, req)
}

// 查看好友申请 
func (*app) ViewAddFriend(ctx context.Context, userid int64, req *emptypb.Empty) (*pb.ViewAddFriendResp, error) {
	list, err := Service.ViewAddFriend(ctx, userid)
	return &pb.ViewAddFriendResp{
		Friends: list,
	}, err
}

// 同意添加好友
func (*app) AgreeFriend(ctx context.Context, userid int64, req *pb.AgreeFriendReq) error {
	return Service.AgreeFriend(ctx, userid, req)
}

// 设置好友信息
func (*app) SetFriend(ctx context.Context, userid int64, req *pb.SetFriendReq) (*pb.SetFriendResp, error) {
	return Service.SetFriend(ctx, userid, req)
}

// 获取所有好友
func (*app) GetAllFriends(ctx context.Context, userid int64, req *emptypb.Empty) (*pb.GetAllFriendResp, error) {
	pbFriends, err := Service.GetAllFriends(ctx, userid)
	return &pb.GetAllFriendResp{
		Friends: pbFriends,
	}, err
}
