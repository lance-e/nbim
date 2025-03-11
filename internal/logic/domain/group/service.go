package group

import (
	"context"
	"fmt"
	"nbim/internal/logic/domain/user"
	"nbim/pkg/protocol/pb"
)

type service struct{}

var Service = new(service)

func (*service) CreateGroup(ctx context.Context, ownerId int64, groupName string, avatarUrl string, introduction string, extra string, membersId []int64) (int64, error) {
	//新建group对象
	group := NewGroup(ownerId, groupName, avatarUrl, introduction, int64(len(membersId)+1), extra, membersId)
	//保存到mysql和redis
	err := group.SaveGroupIntoDB()
	if err != nil {
		return 0, err
	}
	return group.Id, nil
}

func (*service) UpdateGroup(ctx context.Context, userId int64, groupId int64, name string, avatarUrl string, introduction string, extra string) error {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return err
	}
	group.Update(name, avatarUrl, introduction, extra)
	err = group.SaveGroupIntoDB()
	if err != nil {
		return err
	}
	//TODO:PUSH to user

	return nil
}

// 获取群基本信息
func (*service) GetGroupInfo(ctx context.Context, groupId int64) (*GroupInfo, error) {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return nil, err
	}
	return group.GroupInfo, nil
}

func (*service) GetAllGroupUserJoined(ctx context.Context, userId int64) ([]GroupInfo, error) {
	return Dao.ListAllGroupUserJoined(userId)
}

// 获取群全部成员
func (*service) GetGroupAllMember(ctx context.Context, groupId int64) ([]*GroupMember, error) {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return nil, err
	}
	return group.GroupAllMember, nil
}

func (*service) AddGroupMember(ctx context.Context, userId int64, groupId int64, userIds []int64) ([]int64, error) {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return nil, err
	}
	exist, added, err := group.AddMember(userIds)
	if err != nil {
		return nil, err
	}

	err = group.SaveGroupIntoDB()
	if err != nil {
		return nil, err
	}

	//TODO:push
	fmt.Printf("added: %v\n", added)

	return exist, nil
}

func (*service) UpdateGroupMember(ctx context.Context, userId int64, groupId int64, memberType int, remarks string, extra string) error {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return err
	}
	err = group.UpdateMember(userId, memberType, remarks, extra)
	if err != nil {
		return err
	}
	err = group.SaveGroupIntoDB()
	if err != nil {
		return err
	}
	return nil
}

func (*service) DeleteGroupMember(ctx context.Context, admin int64, groupId int64, userId int64) error {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return err
	}

	err = group.DeleteMember(userId)
	if err != nil {
		return nil
	}

	err = group.SaveGroupIntoDB()
	if err != nil {
		return nil
	}

	//TODO:push

	return nil
}

func (*service) GetGroupMembers(ctx context.Context, groupId int64, userId int64) ([]*pb.GroupMember, error) {
	group, err := GetGroupFromDB(groupId)
	if err != nil {
		return nil, err
	}

	//向业务服务器发送获取批量用户信息的rpc请求, 以获取更多成员信息
	userIds := make(map[int64]int32, len(group.GroupAllMember))
	for _, member := range group.GroupAllMember {
		userIds[member.UserId] = 0
	}
	resp, err := user.App.GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	//构建pb.GroupMember列表
	pbAllMembers := make([]*pb.GroupMember, len(group.GroupAllMember))
	for i := range group.GroupAllMember {

		pbMember := &pb.GroupMember{
			UserId:     group.GroupAllMember[i].UserId,
			Remarks:    group.GroupAllMember[i].Remarks,
			Extra:      group.GroupAllMember[i].Extra,
			MemberType: pb.MemberType(group.GroupAllMember[i].MemberType),
		}
		if moreInfo, ok := resp.Users[group.GroupAllMember[i].UserId]; ok {
			pbMember.AvatarUrl = moreInfo.AvatarUrl
			pbMember.Nickname = moreInfo.Nickname
			pbMember.Sex = moreInfo.Sex
			pbMember.UserExtra = moreInfo.Extra
		}

		pbAllMembers[i] = pbMember
	}
	return pbAllMembers, nil
}
