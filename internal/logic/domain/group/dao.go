package group

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"

	"github.com/jinzhu/gorm"
)

type dao struct{}

var Dao = new(dao)

// GetInfo:获取群组信息
func (*dao) GetInfo(groupId int64) (*GroupInfo, error) {
	group := GroupInfo{Id: groupId}
	err := db.DB.First(&group).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &group, err
}

// SaveInfo:保存群组信息
func (*dao) SaveInfo(group *GroupInfo) error {
	return gerror.WrapError(db.DB.Save(&group).Error)
}

// ListAllGroupUserJoined:列举用户加入的所有群组信息
func (*dao) ListAllGroupUserJoined(userId int64) ([]GroupInfo, error) {

	infos := []GroupInfo{}
	err := db.DB.Select("g.id,g.name,g.avatar_url,g.introduction,g.user_num,g.extra,g.create_time,g.update_time").Table("group_member m").Joins("join group_info g on m.group_id = g.id").Where("m.user_id = ?", userId).Find(&infos).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return infos, nil
}

// ListAllMember:获取群组所有成员信息
func (*dao) ListAllMember(groupId int64) ([]*GroupMember, error) {
	members := []*GroupMember{}
	err := db.DB.Find(&members, "group_id = ?", groupId).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return members, nil
}

// GetMember:获取群组指定成员信息,不存在就返回nil
func (*dao) GetMember(groupId int64, userId int64) (*GroupMember, error) {
	member := GroupMember{}
	err := db.DB.First(&member, "group_id = ? and user_id = ?", groupId, userId).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &member, nil
}

// GetMember:保存成员信息
func (*dao) SaveMember(member *GroupMember) error {
	return gerror.WrapError(db.DB.Save(&member).Error)
}

// DeleteMember:删除成员信息
func (*dao) DeleteMember(groupId int64, userId int64) error {
	return gerror.WrapError(db.DB.Exec("delete from group_member where group_id = ? and user_id = ?", groupId, userId).Error)
}
