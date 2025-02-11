package group_test

import (
	"fmt"
	"nbim/internal/logic/domain/group"
	"testing"
	"time"
)

func TestSet(t *testing.T) {

	group1 := group.Group{

		GroupInfo: &group.GroupInfo{
			Id:           1,
			Name:         "vim交流群",
			AvatarUrl:    "https://avatar1",
			Introduction: "vimvimvim",
			UserNum:      0,
			Extra:        "",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		},
		GroupAllMember: make([]group.GroupMember, 0),
	}
	group.Cathe.Set(&group1)
}

func TestGet(t *testing.T) {
	group, err := group.Cathe.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if group != nil {
		fmt.Printf("group1 : id-%d ,name-%s, Introduction-%s\n", group.GroupInfo.Id, group.GroupInfo.Name, group.GroupInfo.Introduction)
	} else {
		fmt.Printf("group is nil\n")
	}
}
