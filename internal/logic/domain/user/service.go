package user

import (
	"errors"
	"time"
)

type service struct{}

var Service = new(service)

func (*service) SignIn(username string, password string, device_id int64) (bool, int64, error) {
	user, err := Dao.GetByUsername(username)
	if err != nil {
		return false, -1, err
	}
	//注册
	if user == nil {
		err = Dao.Save(&User{
			Username:   username,
			Password:   password,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			return true, -1, err
		}
		user, err := Dao.GetByUsername(username)
		if err != nil {
			return true, -1, err
		}
		return true, user.Id, nil
	}
	//正常验证密码
	if user.Password != password {
		return false, -1, errors.New("password wrong")
	} else {
		return false, user.Id, nil
	}
}

func (*service) GetUser(user_id int64) (*User, error) {
	return Dao.GetById(user_id)
}
func (*service) GetUsers(user_ids []int64) ([]User, error) {
	return Dao.GetByIds(user_ids)
}

func (*service) UpdateUser(userId int64, nickname string, sex int32, avatarUrl string, extra string) error {
	user, err := Dao.GetById(userId)
	if err != nil || user == nil {
		return err
	}
	return Dao.Save(&User{
		Id:         user.Id,
		Username:   user.Username,
		Password:   user.Password,
		Nickname:   nickname,
		Sex:        sex,
		AvatarUrl:  avatarUrl,
		Extra:      extra,
		CreateTime: user.CreateTime,
		UpdateTime: time.Now(),
	})
}

func (*service) SearchUser(key string) ([]User, error) {
	return Dao.Search(key)
}
