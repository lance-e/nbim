package room

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type app struct{}

var App = new(app)

// 订阅房间
func (*app) SubscribeRoom(ctx *context.Context, req *pb.SubscribeRoomReq) (*emptypb.Empty, error) {

}
