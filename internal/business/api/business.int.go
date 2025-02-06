package api

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BusinessIntServer struct {
	pb.UnsafeBusinessIntServer
}

// 权限校验
func (s *BusinessIntServer) Auth(context.Context, *pb.AuthReq) (*emptypb.Empty, error) {

	return new(emptypb.Empty), nil
}

// 获取用户信息
func (s *BusinessIntServer) GetUser(context.Context, *pb.GetUserReq) (*pb.GetUserResp, error) {

	return &pb.GetUserResp{}, nil
}

// 批量获取用户信息
func (s *BusinessIntServer) GetUsers(context.Context, *pb.GetUsersReq) (*pb.GetUsersResp, error) {

	return &pb.GetUsersResp{}, nil
}
