package gateway

import (
	"context"
	"nbim/pkg/protocol/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GatewayServer struct {
	pb.UnsafeGatewayServer
}

// gateway server: 发送下行消息
func (g *GatewayServer) SendDownlinkMessage(ctx context.Context, req *pb.SendDownlinkMessageRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
