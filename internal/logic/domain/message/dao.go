package message

import (
	"nbim/pkg/db"
	"nbim/pkg/gerror"
	"time"
)

type dao struct{}

var Dao = new(dao)

//-------------------------all message -------------

// 批量保存
func (*dao) Save(msg *Messages) error {
	err := db.DB.Create(msg).Error
	if err != nil {
		return gerror.WrapError(err)
	}
	return nil
}

// 批量获取消息
func (*dao) GetMany(seqs []int64) ([]*Messages, error) {
	var msgs []*Messages
	err := db.DB.Order("send_time").Find(&msgs, seqs).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return msgs, nil
}

// 软删除
func (*dao) DeleteMany(seqs []int64) error {
	err := db.DB.Table("messages").Where("seq IN ?", seqs).Updates(map[string]interface{}{
		"is_deleted":  1,
		"update_time": time.Now(),
	}).Error
	if err != nil {
		return gerror.WrapError(err)
	}
	return nil

}

//------------------------user message ----------

// 批量保存
func (*dao) SaveUserMsg(msgs []*Messages) error {
	u1 := make([]UserMessages, 2*len(msgs))
	for _, msg := range msgs {
		u1 = append(u1, UserMessages{
			UserId:      int64(msg.SessionId),
			Seq:         msg.Seq,
			ReceiveTime: time.Now().UnixMilli(),
			Status:      0, //接收者未读
			IsDeleted:   0,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		})
		u1 = append(u1, UserMessages{
			UserId:      msg.SenderId,
			Seq:         msg.Seq,
			ReceiveTime: time.Now().UnixMilli(),
			Status:      1, //发送者肯定已读
			IsDeleted:   0,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		})
	}
	err := db.DB.Create(&u1).Error
	if err != nil {
		return gerror.WrapError(err)
	}
	return nil
}

// TODO:先暂时拉取全部未读消息，可优化为分页查询，未读消息过多就每次少量拉取最新的几十条消息(status 结合 receiveTime)
func (*dao) GetUserMsgIds(ownerId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := db.DB.Model(&UserMessages{}).Where("status = ?", 0).Pluck("seq", ids).Error
	if err != nil {
		return nil, gerror.WrapError(err)
	}
	return ids, nil
}

//TODO:历史消息
