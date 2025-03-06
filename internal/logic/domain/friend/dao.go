package friend

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"

	"github.com/jinzhu/gorm"
)

type dao struct{}

var Dao = new(dao)

// Get:获取用户的指定好友信息
func (*dao) Get(userId int64, friendId int64) (*Friend, error) {
	friend := Friend{}
	err := db.DB.First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &friend, nil
}

// Save:保存好友信息 
func (*dao) Save(friend *Friend) error {
	return gerror.WrapError(db.DB.Save(&friend).Error)
}

// List:获取符合指定状态的好友关系列表(好友列表或者申请列表)
func (*dao) List(userId int64, status int) ([]Friend, error) {
	friends := []Friend{}
	var err error
	if status == 1 {
		//查看好友列表
		err = db.DB.Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	} else {
		//查看好友申请
		err = db.DB.Where("friend_id= ? and status = ?", userId, status).Find(&friends).Error

	}
	return friends, gerror.WrapError(err)
}
