package group

import (
	"errors"
	"fmt"
	"nbim/pkg/protocol/pb"
	"time"
)

const (
	MemberChangeUpdate = 1
	MemberChangeDelete = 2
)

// 群组全部信息
type Group struct {
	*GroupInfo                    //基本信息
	GroupAllMember []*GroupMember //成员
}

// GroupInfo:群组信息
type GroupInfo struct {
	Id           int64
	Name         string
	AvatarUrl    string
	Introduction string
	UserNum      int64
	Extra        string
	CreateTime   time.Time
	UpdateTime   time.Time
}

// 群组成员信息
type GroupMember struct {
	Id         int64
	GroupId    int64
	UserId     int64
	MemberType int
	Remarks    string
	Extra      string
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
	ChangeType int `gorm:"-"` //仅用于逻辑处理
}

// NewGroup:创建群组实例，包含群组信息和群组成员信息
func NewGroup(ownerId int64, name string, avatarUrl string, introduction string, userNum int64, extra string, membersIds []int64) *Group {
	//(数据库还未分配group_id)
	group := &Group{
		GroupInfo:      NewGroupInfo(name, avatarUrl, introduction, userNum, extra),
		GroupAllMember: make([]*GroupMember, 0),
	}
	//将建群用户特殊处理为管理员
	group.GroupAllMember = append(group.GroupAllMember, NewGroupMember(0, ownerId, int(pb.MemberType_GMT_ADMIN)))
	//将其他成员添加进群组，并设置为普通成员
	for i := range membersIds {
		group.GroupAllMember = append(group.GroupAllMember, NewGroupMember(0, membersIds[i], int(pb.MemberType_GMT_MEMBER)))
	}
	return group
}

func NewGroupInfo(name string, avatarUrl string, introduction string, userNum int64, extra string) *GroupInfo {
	info := &GroupInfo{
		Name:         name,
		AvatarUrl:    avatarUrl,
		Introduction: introduction,
		UserNum:      userNum,
		Extra:        extra,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	return info
}

func NewGroupMember(groupId int64, userId int64, memberType int) *GroupMember {
	member := &GroupMember{
		GroupId:    groupId,
		UserId:     userId,
		MemberType: memberType,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ChangeType: MemberChangeUpdate,
	}
	return member
}

// GetGroupFromDB:从redis或mysql中获取群组所有信息
func GetGroupFromDB(groupId int64) (*Group, error) {
	//先在缓存中找
	group, err := Cathe.Get(groupId)
	if err != nil {
		return nil, err
	}
	if group != nil {
		return group, nil
	}
	//缓存未找到，去数据库找
	info, err := Dao.GetInfo(groupId)
	if err != nil {
		return nil, err
	}
	members, err := Dao.ListAllMember(groupId)
	if err != nil {
		return nil, err
	}
	group = &Group{
		GroupInfo:      info,
		GroupAllMember: members,
	}
	//从数据库查后，缓存
	err = Cathe.Set(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// SaveGroupIntoDB:保存群组所有信息到数据库，并返回群组group_id
// !!!一次性保存批量更改，减少数据库io
func (group *Group) SaveGroupIntoDB() error {
	groupId := group.GroupInfo.Id
	//先保存群组信息,(若第一次创建会自动分配group_id)
	if err := Dao.SaveInfo(group.GroupInfo); err != nil {
		return err
	}
	//分别将每个成员信息保存到数据库
	for _, member := range group.GroupAllMember {
		member.GroupId = group.GroupInfo.Id //与群组id绑定
		switch member.ChangeType {
		case MemberChangeUpdate:
			if err := Dao.SaveMember(member); err != nil {
				return err
			}
		case MemberChangeDelete:
			fmt.Printf("执行delete\n")
			if err := Dao.DeleteMember(member.GroupId, member.UserId); err != nil {
				return err
			}
		}
	}

	//若groupId最初不为0,则说明群组此前已存在，此次为更新群组信息，之前的群组缓存失效，直接删除
	if groupId != 0 {
		err := Cathe.Del(groupId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (group *Group) Update(name string, avatarUrl string, introduction string, extra string) {
	group.Name = name
	group.AvatarUrl = avatarUrl
	group.Introduction = introduction
	group.Extra = extra
}

func (group *Group) IsMember(userId int64) bool {
	for _, member := range group.GroupAllMember {
		if member.UserId == userId {
			return true
		}
	}
	return false
}

func (group *Group) GetMember(userId int64) *GroupMember {
	for _, member := range group.GroupAllMember {
		if member.UserId == userId {
			return member
		}
	}
	return nil
}

// 群组添加成员，返回之前已加入的成员id，此次加入的成员id,error
func (group *Group) AddMember(userIds []int64) ([]int64, []int64, error) {
	exist := make([]int64, 0)
	added := make([]int64, 0)
	for _, userId := range userIds {
		//已经是成员
		if group.IsMember(userId) {
			exist = append(exist, userId)
			continue
		}
		//非成员，加入群组
		group.GroupAllMember = append(group.GroupAllMember, &GroupMember{
			GroupId:    group.Id,
			UserId:     userId,
			MemberType: int(pb.MemberType_GMT_MEMBER),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			ChangeType: MemberChangeUpdate,
		})
		added = append(added, userId)
	}
	return exist, added, nil
}

func (group *Group) UpdateMember(userId int64, memberType int, remarks string, extra string) error {
	member := group.GetMember(userId)
	if member == nil {
		return errors.New("update member failed")
	}

	member.MemberType = memberType
	member.Remarks = remarks
	member.Extra = extra
	member.UpdateTime = time.Now()
	member.ChangeType = MemberChangeUpdate

	return nil
}

func (group *Group) DeleteMember(userId int64) error {
	member := group.GetMember(userId)
	if member == nil {
		return errors.New("delete member failed")
	}

	member.ChangeType = MemberChangeDelete
	return nil
}
