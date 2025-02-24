package gateway

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ConnectionIntServer struct {
	pb.UnsafeConnectionIntServer
}

// 消息转发
func (s *ConnectionIntServer) TransferMessage(context.Context, *pb.TransferMessageReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
