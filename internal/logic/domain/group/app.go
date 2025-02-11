package group

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 创建群聊
func (*app) CreateGroup(ctx context.Context, userid int64, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	groupId, err := Service.CreateGroup(ctx, userid, req.Name, req.AvatarUrl, req.Introduction, req.Extra, req.MemberIds)
	return &pb.CreateGroupResp{
		GroupId: groupId,
	}, err
}

// 更新群组
func (*app) UpdateGroup(ctx context.Context, userid int64, req *pb.UpdateGroupReq) error {
	return Service.UpdateGroup(ctx, userid, req.GroupId, req.Name, req.AvatarUrl, req.AvatarUrl, req.Extra)
}

// 获取群组信息
func (*app) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	info, err := Service.GetGroupInfo(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	return &pb.GetGroupResp{
		Group: &pb.GroupInfo{
			GroupId:      req.GroupId,
			AvatarUrl:    info.AvatarUrl,
			Name:         info.Name,
			Introduction: info.Introduction,
			UserNum:      info.UserNum,
			Extra:        info.Extra,
			CreateTime:   info.CreateTime.UnixMilli(),
			UpdateTime:   info.UpdateTime.UnixMilli(),
		},
	}, nil
}

// 获取用户加入的所有群组信息
func (*app) GetAllGroup(ctx context.Context, userid int64, req *emptypb.Empty) (*pb.GetAllGroupResp, error) {
	infos, err := Service.GetAllGroupUserJoined(ctx, userid)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetAllGroupResp{}
	for _, info := range infos {
		resp.Groups = append(resp.Groups, &pb.GroupInfo{GroupId: info.Id,
			AvatarUrl:    info.AvatarUrl,
			Name:         info.Name,
			Introduction: info.Introduction,
			UserNum:      info.UserNum,
			Extra:        info.Extra,
			CreateTime:   info.CreateTime.UnixMilli(),
			UpdateTime:   info.UpdateTime.UnixMilli(),
		})
	}
	return resp, nil
}

// 添加群组成员
func (*app) AddGroupMember(ctx context.Context, userid int64, req *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	exist, err := Service.AddGroupMember(ctx, userid, req.GroupId, req.UserIds)
	if err != nil {
		return nil, err
	}
	return &pb.AddGroupMemberResp{UserIds: exist}, nil
}

// 更新群组成员信息
func (*app) UpdateGroupMember(ctx context.Context, req *pb.UpdateGroupMemberReq) error {
	return Service.UpdateGroupMember(ctx, req.UserId, req.GroupId, int(req.MemberType), req.Remarks, req.Extra)
}

// 删除群组成员
func (*app) DeleteGroupMember(ctx context.Context, admin int64, req *pb.DeleteGroupMemberReq) error {
	return Service.DeleteGroupMember(ctx, admin, req.GroupId, req.UserId)
}

// 获取群组成员
func (*app) GetGroupMember(ctx context.Context, req *pb.GetGroupMemberReq) (*pb.GetGroupMemberResp, error) {
	members, err := Service.GetGroupMembers(ctx, req.GroupId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetGroupMemberResp{
		Members: members,
	}, nil
}
