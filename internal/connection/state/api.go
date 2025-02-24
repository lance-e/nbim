package state

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StateServer struct {
	pb.UnsafeStateServer
}

// state server
func (g *StateServer) ReceiveUplinkMessage(ctx context.Context, req *pb.ReceiveUplinkMessageRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
