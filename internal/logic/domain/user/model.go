package user

import "time"

// 用户信息
type User struct {
	Id         int64
	Username   string
	Password   string
	Nickname   string
	Sex        int32
	AvatarUrl  string
	Extra      string
	CreateTime time.Time
	UpdateTime time.Time
}
