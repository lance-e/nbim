package proxy

import (
	"context"
	"nbim/pkg/protocol/pb"
)

var DevcieProxy deviceProxy

type deviceProxy interface {
	ListAllOnlineDeviceByUserId(ctx context.Context, userId int64) ([]*pb.Device, error)
}
