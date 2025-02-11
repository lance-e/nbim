package friend

import (
	"context"
	"fmt"
	"nbim/pkg/gerror"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"time"
)

type service struct{}

var Service = new(service)

func (*service) AddFriend(ctx context.Context, userId int64, req *pb.AddFriendReq) error {
	friend, err := Dao.Get(userId, req.FriendId)
	if err != nil {
		return err
	}
	//好友已存在
	if friend != nil {
		if friend.Status == FriendStatusApply {
			return nil
		} else if friend.Status == FriendStatusAgree {
			return gerror.ErrAlreadyIsFriend
		}
	}

	//创建好友关系(申请状态)
	err = Dao.Save(&Friend{
		UserId:     userId,
		FriendId:   req.FriendId,
		Remarks:    req.Remarks,
		Status:     FriendStatusApply,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		return err
	}

	//调用业务服务器接口,获取用户信息
	resp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}

	//TODO
	fmt.Printf("resp: %v\n", resp)

	return nil
}

func (*service) AgreeFriend(ctx context.Context, userId int64, req *pb.AgreeFriendReq) error {
	friend, err := Dao.Get(req.FriendId, userId)
	if err != nil {
		return err
	}
	if friend == nil {
		return gerror.ErrBadRequest
	}
	if friend.Status == FriendStatusAgree {
		return nil
	}

	//保存申请方作为起点，被申请方作为终点的好友关系
	friend.Status = FriendStatusAgree
	err = Dao.Save(friend)
	if err != nil {
		return err
	}

	//保存申请方作为重点，被申请方作为起点的好友关系
	err = Dao.Save(&Friend{
		UserId:     userId,
		FriendId:   req.FriendId,
		Remarks:    friend.Remarks,
		Status:     FriendStatusAgree,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		return err
	}

	//调用业务服务器接口,获取用户信息
	resp, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}

	//TODO
	fmt.Printf("resp: %v\n", resp)

	return nil
}

func (*service) SetFriend(ctx context.Context, userId int64, req *pb.SetFriendReq) (*pb.SetFriendResp, error) {
	friend, err := Dao.Get(userId, req.FriendId)
	if err != nil {
		return nil, err
	}
	if friend == nil {
		return nil, nil
	}
	//更改信息
	friend.Remarks = req.Remarks
	friend.Extra = req.Extra
	friend.UpdateTime = time.Now()

	err = Dao.Save(friend)
	if err != nil {
		return nil, err
	}
	return &pb.SetFriendResp{
		FriendId: req.FriendId,
		Remarks:  req.Remarks,
		Extra:    req.Extra,
	}, nil
}

func (*service) GetAllFriends(ctx context.Context, userId int64) ([]*pb.Friend, error) {
	friends, err := Dao.List(userId, FriendStatusAgree)
	if err != nil {
		return nil, err
	}

	//只使用到了key
	userIds := make(map[int64]int32, len(friends))
	for _, friend := range friends {
		userIds[friend.FriendId] = 0
	}

	//通过调用业务服务器获取更多的用户信息
	resp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{
		UserIds: userIds,
	})
	if err != nil {
		return nil, err
	}

	//构造返回的pb.Friend
	pbFriends := make([]*pb.Friend, len(friends))
	for i := range friends {
		friend := pb.Friend{
			FriendId: friends[i].FriendId,
			Remarks:  friends[i].Remarks,
			Extra:    friends[i].Extra,
		}
		if res, ok := resp.Users[friends[i].FriendId]; ok {
			friend.Nickname = res.Nickname
			friend.Sex = res.Sex
			friend.AvatarUrl = res.AvatarUrl
			friend.UserExtra = res.Extra
		}
		pbFriends[i] = &friend
	}
	return pbFriends, nil

}
