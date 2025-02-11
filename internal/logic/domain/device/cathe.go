package device

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type cathe struct{}

var (
	//key:用户 value:在线设备列表
	UserDeviceRedisKey    = "user_device:"
	UserDeviceRedisExpire = 2 * time.Hour
)

var Cathe = new(cathe)

// Get:获取指定用户的全部在线设备
func (c *cathe) Get(userId int64) ([]Device, error) {
	var devices []Device = make([]Device, 0)
	err := db.GetRedisByJson(UserDeviceRedisKey+strconv.FormatInt(userId, 10), &devices)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, gerror.WrapError(err)
	}
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return devices, nil
}

// Set:设置指定用户的全部在线设备
func (c *cathe) Set(userId int64, devices []Device) error {
	err := db.SetRedisByJson(UserDeviceRedisKey+strconv.FormatInt(userId, 10), devices, UserDeviceRedisExpire)
	return gerror.WrapError(err)
}

// Del:删除指定用户的在线设备列表
func (c *cathe) Del(userId int64) error {
	_, err := db.RedisCli.Del(UserDeviceRedisKey + strconv.FormatInt(userId, 10)).Result()
	return gerror.WrapError(err)
}
