package user

import (
	"context"
	"errors"
	"nbim/pkg/protocol/pb"
)

type app struct {
}

var App = new(app)

func (*app) SignIn(c context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	is_new, user_id, err := Service.SignIn(req.GetUsername(), req.GetPassword(), req.GetDeviceId())
	return &pb.SignInResp{
		IsNew:  is_new,
		UserId: user_id,
	}, err
}

func (*app) GetUser(c context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := Service.GetUser(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return &pb.GetUserResp{
		User: &pb.User{
			UserId:    user.Id,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Sex:       user.Sex,
			AvatarUrl: user.AvatarUrl,
			Extra:     user.Extra,
		},
	}, nil
}

func (*app) GetUsers(c context.Context, req *pb.GetUsersReq) (*pb.GetUsersResp, error) {
	ids := make([]int64, len(req.UserIds))
	for k := range req.UserIds {
		ids = append(ids, k)
	}
	users, err := Service.GetUsers(ids)
	if err != nil {
		return nil, err
	}
	pbUser := make(map[int64]*pb.User, len(users))
	for i := range users {
		pbUser[users[i].Id] = &pb.User{
			UserId:    users[i].Id,
			Username:  users[i].Username,
			Nickname:  users[i].Nickname,
			Sex:       users[i].Sex,
			AvatarUrl: users[i].AvatarUrl,
			Extra:     users[i].Extra,
		}
	}
	return &pb.GetUsersResp{
		Users: pbUser,
	}, nil
}

func (*app) UpdateUser(c context.Context, userId int64, req *pb.UpdateUserReq) error {
	return Service.UpdateUser(userId, req.Nickname, req.Sex, req.AvatarUrl, req.Extra)
}

func (*app) SearchUser(c context.Context, req *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	users, err := Service.SearchUser(req.Key)
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*pb.User, len(users))
	for i := range users {
		pbUsers[i] = &pb.User{
			UserId:    users[i].Id,
			Username:  users[i].Username,
			Nickname:  users[i].Nickname,
			Sex:       users[i].Sex,
			AvatarUrl: users[i].AvatarUrl,
			Extra:     users[i].Extra,
		}
	}
	return &pb.SearchUserResp{Users: pbUsers}, nil
}
