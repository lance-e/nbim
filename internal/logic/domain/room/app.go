package room

import (
	"context"
	"nbim/pkg/protocol/pb"
)

type app struct{}

var App = new(app)

// 订阅房间
func (*app) SubscribeRoom(ctx *context.Context, req *pb.SubscribeRoomReq) error {

}
