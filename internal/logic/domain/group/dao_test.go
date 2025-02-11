package group_test

import (
	"fmt"
	"nbim/internal/logic/domain/group"
	"nbim/pkg/protocol/pb"
	"testing"
	"time"
)

func TestGroupInfoSave(t *testing.T) {
	info1 := group.GroupInfo{
		Name:         "vim交流群",
		AvatarUrl:    "https://avatar1",
		Introduction: "vimvimvim",
		UserNum:      0,
		Extra:        "",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	info2 := group.GroupInfo{
		Name:         "go交流群",
		AvatarUrl:    "https://avatar2",
		Introduction: "go是世界上最好的语言",
		UserNum:      0,
		Extra:        "",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	info3 := group.GroupInfo{
		Name:         "27实习群",
		AvatarUrl:    "https://avatar3",
		Introduction: "日常实习",
		UserNum:      0,
		Extra:        "",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	group.Dao.SaveInfo(&info1)
	group.Dao.SaveInfo(&info2)
	group.Dao.SaveInfo(&info3)
}

func TestGroupInfoGet(t *testing.T) {
	info1, err := group.Dao.GetInfo(1)
	if err != nil {
		t.Fatal(err)
	}
	info2, err := group.Dao.GetInfo(3)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("GroupInfo1 : id-%d, name-%s, Introduction:%s\n", info1.Id, info1.Name, info1.Introduction)
	fmt.Printf("GroupInfo2 : id-%d, name-%s, Introduction:%s\n", info2.Id, info2.Name, info2.Introduction)
}

func TestSaveMember(t *testing.T) {
	member11 := group.GroupMember{
		GroupId:    1,
		UserId:     1,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member21 := group.GroupMember{
		GroupId:    2,
		UserId:     1,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member31 := group.GroupMember{
		GroupId:    3,
		UserId:     1,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		Remarks:    "我是在群3的用户1",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member12 := group.GroupMember{
		GroupId:    1,
		UserId:     2,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member13 := group.GroupMember{
		GroupId:    1,
		UserId:     3,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		Remarks:    "我是在群1的用户3",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member14 := group.GroupMember{
		GroupId:    1,
		UserId:     4,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	member15 := group.GroupMember{
		GroupId:    1,
		UserId:     5,
		MemberType: int(pb.MemberType_GMT_MEMBER),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	group.Dao.SaveMember(&member11)
	group.Dao.SaveMember(&member21)
	group.Dao.SaveMember(&member31)
	group.Dao.SaveMember(&member12)
	group.Dao.SaveMember(&member13)
	group.Dao.SaveMember(&member14)
	group.Dao.SaveMember(&member15)
}

func TestGetMember(t *testing.T) {
	member1, err := group.Dao.GetMember(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	member2, err := group.Dao.GetMember(1, 3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("member1 Remarks:%s\n", member1.Remarks)
	fmt.Printf("member2 Remarks:%s\n", member2.Remarks)
}

func TestListAllGroupUserJoined(t *testing.T) {
	infos, err := group.Dao.ListAllGroupUserJoined(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range infos {
		fmt.Printf("info: group_id-%d, name-%s, introduction-%s\n", info.Id, info.Name, info.Introduction)
	}
}

func TestListAllGroupMember(t *testing.T) {
	members, err := group.Dao.ListAllMember(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, member := range members {
		fmt.Printf("member: Id-%d\n", member.UserId)
	}
}

func TestDeleteMember(t *testing.T) {
	err := group.Dao.DeleteMember(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	group, err := group.Dao.GetMember(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	if group != nil {
		t.Fatal("group3, user1 not deleted\n")
	}
}
