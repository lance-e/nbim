package api

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BusinessExtServer struct {
	pb.UnsafeBusinessExtServer
}

// 登陆
func (s *BusinessExtServer) SignIn(context.Context, *pb.SignInReq) (*pb.SignInResp, error) {
	return &pb.SignInResp{}, nil
}

// 获取用户信息
func (s *BusinessExtServer) GetUser(context.Context, *pb.GetUserReq) (*pb.GetUserResp, error) {
	return &pb.GetUserResp{}, nil
}

// 更新用户信息
func (s *BusinessExtServer) UpdateUser(context.Context, *pb.UpdateUserReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// 搜索用户
func (s *BusinessExtServer) SearchUser(context.Context, *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	return &pb.SearchUserResp{}, nil
}
