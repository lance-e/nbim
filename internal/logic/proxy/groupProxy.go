package proxy

import (
	"context"
	"nbim/internal/logic/domain/group"
)

var GroupProxy groupProxy

type groupProxy interface {
	GetGroupAllMember(ctx context.Context, groupId int64) ([]*group.GroupMember, error)
}
