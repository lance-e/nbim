package api

import (
	"context"
	"nbim/internal/logic/domain/device"
	"nbim/internal/logic/domain/friend"
	"nbim/internal/logic/domain/group"
	"nbim/internal/logic/domain/message"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicExtServer struct {
	pb.UnsafeLogicExtServer
}

// 注册设备
func (s *LogicExtServer) RegisterDevice(ctx context.Context, req *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {
	return device.App.RegisterDevice(ctx, req)
}

// 推送信息到房间
func (s *LogicExtServer) PushRoom(ctx context.Context, req *pb.PushRoomReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, message.App.PushRoom(ctx, req)
}

// 发送好友消息
func (s *LogicExtServer) SendMessageToFriend(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	userid, deviceid, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return message.App.SendMessageToFriend(ctx, userid, deviceid, req)
}

// 添加好友
func (s *LogicExtServer) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*emptypb.Empty, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, friend.App.AddFriend(ctx, userid, req)
}

// 同意添加好友
func (s *LogicExtServer) AgreeFriend(ctx context.Context, req *pb.AgreeFriendReq) (*emptypb.Empty, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, friend.App.AgreeFriend(ctx, userid, req)
}

// 设置好友信息
func (s *LogicExtServer) SetFriend(ctx context.Context, req *pb.SetFriendReq) (*pb.SetFriendResp, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return friend.App.SetFriend(ctx, userid, req)
}

// 获取所有好友
func (s *LogicExtServer) GetAllFriends(ctx context.Context, req *emptypb.Empty) (*pb.GetAllFriendResp, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return friend.App.GetAllFriends(ctx, userid, req)
}

// 发送群组信息
func (s *LogicExtServer) SendMessageToGroup(ctx context.Context, req *pb.SendMessageReq) (*pb.SendMessageResp, error) {
	userid, deviceid, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return message.App.SendMessageToGroup(ctx, userid, deviceid, req)
}

// 创建群聊
func (s *LogicExtServer) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return group.App.CreateGroup(ctx, userid, req)
}

// 更新群组
func (s *LogicExtServer) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) (*emptypb.Empty, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, group.App.UpdateGroup(ctx, userid, req)
}

// 获取群组信息
func (s *LogicExtServer) GetGroup(ctx context.Context, req *pb.GetGroupReq) (*pb.GetGroupResp, error) {
	return group.App.GetGroup(ctx, req)
}

// 获取用户加入的所有群组信息
func (s *LogicExtServer) GetAllGroup(ctx context.Context, req *emptypb.Empty) (*pb.GetAllGroupResp, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return group.App.GetAllGroup(ctx, userid, req)
}

// 添加群组成员
func (s *LogicExtServer) AddGroupMember(ctx context.Context, req *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return group.App.AddGroupMember(ctx, userid, req)
}

// 更新群组成员信息
func (s *LogicExtServer) UpdateGroupMember(ctx context.Context, req *pb.UpdateGroupMemberReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, group.App.UpdateGroupMember(ctx, req)
}

// 删除群组成员
func (s *LogicExtServer) DeleteGroupMember(ctx context.Context, req *pb.DeleteGroupMemberReq) (*emptypb.Empty, error) {
	userid, _, err := rpc.GetCtxUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, group.App.DeleteGroupMember(ctx, userid, req)
}

// 获取群组成员
func (s *LogicExtServer) GetGroupMember(ctx context.Context, req *pb.GetGroupMemberReq) (*pb.GetGroupMemberResp, error) {
	return group.App.GetGroupMember(ctx, req)
}
