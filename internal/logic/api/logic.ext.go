package api

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type LogicExtServer struct {
	pb.UnsafeLogicExtServer
}

// 注册设备
func (s *LogicExtServer) RegisterDevice(context.Context, *pb.RegisterDeviceReq) (*pb.RegisterDeviceResp, error) {

	return &pb.RegisterDeviceResp{}, nil
}

// 推送信息到房间
func (s *LogicExtServer) PushRoom(context.Context, *pb.PushRoomReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 发送好友消息
func (s *LogicExtServer) SendMessageToFriend(context.Context, *pb.SendMessageReq) (*pb.SendMessageResp, error) {

	return &pb.SendMessageResp{}, nil
}

// 添加好友
func (s *LogicExtServer) AddFriend(context.Context, *pb.AddFriendReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 同意添加好友
func (s *LogicExtServer) AgreeFriend(context.Context, *pb.AgreeFriendReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 设置好友信息
func (s *LogicExtServer) SetFriend(context.Context, *pb.SetFriendReq) (*pb.SetFriendResp, error) {

	return &pb.SetFriendResp{}, nil
}

// 获取所有好友
func (s *LogicExtServer) GetAllFriends(context.Context, *emptypb.Empty) (*pb.GetAllFriendResp, error) {

	return &pb.GetAllFriendResp{}, nil
}

// 发送群组信息
func (s *LogicExtServer) SendMessageToGroup(context.Context, *pb.SendMessageReq) (*pb.SendMessageResp, error) {

	return &pb.SendMessageResp{}, nil
}

// 创建群聊
func (s *LogicExtServer) CreateGroup(context.Context, *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {

	return &pb.CreateGroupResp{}, nil
}

// 更新群组
func (s *LogicExtServer) UpdateGroup(context.Context, *pb.UpdateGroupReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 获取群组信息
func (s *LogicExtServer) GetGroup(context.Context, *pb.GetGroupReq) (*pb.GetGroupResp, error) {

	return &pb.GetGroupResp{}, nil
}

// 获取用户加入的所有群组信息
func (s *LogicExtServer) GetAllGroup(context.Context, *emptypb.Empty) (*pb.GetAllGroupResp, error) {

	return &pb.GetAllGroupResp{}, nil
}

// 添加群组成员
func (s *LogicExtServer) AddGroupMember(context.Context, *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {

	return &pb.AddGroupMemberResp{}, nil
}

// 更新群组成员信息
func (s *LogicExtServer) UpdateGroupMember(context.Context, *pb.UpdateGroupMemberReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 删除群组成员
func (s *LogicExtServer) DeleteGroupMember(context.Context, *pb.DeleteGroupMemberReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 获取群组成员
func (s *LogicExtServer) GetGroupMember(context.Context, *pb.GetGroupMemberReq) (*pb.GetGroupMemberResp, error) {

	return &pb.GetGroupMemberResp{}, nil
}
