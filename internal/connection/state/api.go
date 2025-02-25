package state

import (
	"context"
	"fmt"
	"nbim/pkg/logger"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StateServer struct {
	pb.UnsafeStateServer
}

// state server：传递上行消息，由长连接网关服务器调用
func (g *StateServer) ReceiveUplinkMessage(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	logger.Logger.Debug("ReceiveUplinkMessage Rpc")
	fmt.Printf("%d\n", req.GetConnId())
	rpc.GetGatewayClient().SendDownlinkMessage(ctx, &pb.GatewayRequest{
		ConnId:  req.GetConnId(),
		Message: req.GetMessage(),
	})
	return new(emptypb.Empty), nil
}

// state server：清除连接状态，由长连接网关服务器调用
func (g *StateServer) ClearConnState(ctx context.Context, req *pb.StateRequest) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
