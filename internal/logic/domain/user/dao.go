package user

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"

	"github.com/jinzhu/gorm"
)

type dao struct{}

var Dao = new(dao)

func (*dao) Save(user *User) error {
	err := db.DB.Save(&user).Error
	if err != nil {
		return gerror.WrapError(err)
	}
	return nil
}

func (*dao) GetById(id int64) (*User, error) {
	var user = User{Id: id}
	err := db.DB.First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerror.WrapError(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, nil
}
func (*dao) GetByIds(ids []int64) ([]User, error) {
	var users []User
	err := db.DB.Find(&users, "id in (?)", ids).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}

	return users, nil
}

func (*dao) GetByUsername(username string) (*User, error) {
	var user User
	err := db.DB.First(&user, "username = ?", username).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerror.WrapError(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, nil
}

func (*dao) Search(key string) ([]User, error) {
	var users []User
	key = "%" + key + "%"
	err := db.DB.Where("username like ? or nickname like ?", key, key).Find(&users).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return users, nil
}
