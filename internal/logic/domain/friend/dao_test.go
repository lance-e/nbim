package friend_test

import (
	"nbim/internal/logic/domain/friend"
	"testing"
	"time"
)

func TestFriend_Save(t *testing.T) {
	f2 := friend.Friend{
		UserId:     1,
		FriendId:   2,
		Remarks:    "朋友2",
		Status:     friend.FriendStatusApply,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	f3 := friend.Friend{
		UserId:     1,
		FriendId:   3,
		Remarks:    "朋友3",
		Status:     friend.FriendStatusAgree,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	f4 := friend.Friend{
		UserId:     1,
		FriendId:   4,
		Remarks:    "朋友4",
		Status:     friend.FriendStatusAgree,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	friend.Dao.Save(&f2)
	friend.Dao.Save(&f3)
	friend.Dao.Save(&f4)
}
func TestFriend_Get(t *testing.T) {
	f2, err := friend.Dao.Get(1, 2)
	if err != nil {
		t.Error(err)
	}
	f3, err := friend.Dao.Get(1, 3)
	if err != nil {
		t.Error(err)
	}

	t.Logf("f2:\n userId:%d\tfriendId:%d\tremarks:%s\tstatus:%d\t\n", f2.UserId, f2.FriendId, f2.Remarks, f2.Status)
	t.Logf("f3:\n userId:%d\tfriendId:%d\tremarks:%s\tstatus:%d\t\n", f3.UserId, f3.FriendId, f3.Remarks, f3.Status)
}

func TestFriend_List(t *testing.T) {
	agreeFriends, err := friend.Dao.List(1, 1)
	if err != nil {
		t.Error(err)
	}
	for i := range agreeFriends {
		t.Logf("agreeFriend: friend_id:%d\n", agreeFriends[i].FriendId)
	}
}
