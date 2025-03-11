package message

import (
	"fmt"
	"nbim/pkg/db"
)

type seq struct{}

var Seq = new(seq)

func (*seq) GetGroupMsgSeq(groupId int64) (int64, error) {
	key := fmt.Sprintf(db.GroupSeqIDKey, groupId)
	result, err := db.RedisCli.Exists(key).Result()
	if err != nil {
		return -1, err
	}
	if result == 0 {
		db.RedisCli.Set(key, 0, db.TTL7Day)
	}
	result, err = db.RedisCli.Incr(key).Result()
	if err != nil {
		return -1, err
	}
	return result, nil
}

func (*seq) GetUserMsgSeq() (int64, error) {

}
