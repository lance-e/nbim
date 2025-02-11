package group

import (
	"errors"
	"nbim/pkg/db"
	"nbim/pkg/gerror"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type cathe struct{}

const (
	//key是group_id,value是序列化后的群组全部信息
	GroupRedisKey    = "group:"
	GroupRedisExpire = 24 * time.Hour
)

var Cathe = new(cathe)

func (*cathe) Get(groupId int64) (*Group, error) {
	group := Group{}
	err := db.GetRedisByJson(GroupRedisKey+strconv.FormatInt(groupId, 10), &group)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, gerror.WrapError(err)
	}
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return &group, nil
}

func (*cathe) Set(group *Group) error {
	return gerror.WrapError(db.SetRedisByJson(GroupRedisKey+strconv.FormatInt(group.Id, 10), group, GroupRedisExpire))
}

func (*cathe) Del(groupId int64) error {
	_, err := db.RedisCli.Del(GroupRedisKey + strconv.FormatInt(groupId, 10)).Result()
	return gerror.WrapError(err)
}
