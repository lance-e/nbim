package group

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 创建群聊
func (*app) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {

}

// 更新群组
func (*app) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) (*emptypb.Empty, error) {

}

// 获取群组信息
func (*app) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.GetGroupResp, error) {

}

// 获取用户加入的所有群组信息
func (*app) GetAllGroup(ctx context.Context, req *emptypb.Empty) (*pb.GetAllGroupResp, error) {

}

// 添加群组成员
func (*app) AddGroupMember(ctx context.Context, req *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {

}

// 更新群组成员信息
func (*app) UpdateGroupMember(ctx context.Context, req *pb.UpdateGroupMemberReq) (*emptypb.Empty, error) {

}

// 删除群组成员
func (*app) DeleteGroupMember(ctx context.Context, req *pb.DeleteGroupMemberReq) (*emptypb.Empty, error) {

}

// 获取群组成员
func (*app) GetGroupMember(ctx context.Context, req *pb.GetGroupMemberReq) (*pb.GetGroupMemberResp, error) {
}
