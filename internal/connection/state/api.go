package state

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StateServer struct {
	pb.UnsafeStateServer
}

// state server：传递上行消息，由长连接网关服务器调用
func (g *StateServer) ReceiveUplinkMessage(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

// state server：，由长连接网关服务器调用
func (g *StateServer) CloseConn(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
