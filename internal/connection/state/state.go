package state

import (
	"context"
	"fmt"
	"nbim/pkg/timer"
	"sync"
)

var CS *catheState

type catheState struct {
	MessageId     int64
	ConnIdToState sync.Map
	Server        *StateServer
}

// 连接状态
type connState struct {
	connID         int64
	deviceID       int64
	messageTimer   *timer.TimeWheel //消息定时器
	reconnectTimer *timer.TimeWheel //重连定时器
	heartbeatTimer *timer.TimeWheel //心跳定时器
}

func InitCatheState() {
	CS = &catheState{}
	CS.Server = &StateServer{cmdchannel: make(chan *cmdContext, 2048)}
}

func (c *catheState) ClearConnState(ctx *context.Context, connID int64) {

}

// 登陆
func (c *catheState) connLogin(ctx *context.Context, deviceId int64, connID int64) error {
	fmt.Printf("CS.connLogin\n")
	return nil
}

func (c *catheState) compareAndIncrementClientID(ctx *context.Context, connID int64, clientID int64, sessionID int64) bool {
	fmt.Printf("CS.compareAndIncrementClientID\n")
	return true
}

func (c *catheState) connReconn(ctx *context.Context, oldConnID, newConnID int64) error {
	fmt.Print("CS.connReconn\n")
	return nil
}

func (c *catheState) connAck(ctx *context.Context, connID int64, sessionID int64, msgID int64) {
	fmt.Print("CS.connAck\n")
}

// 重置心跳计时器
func (c *catheState) reSetHeartbeatTimer(connId int64) {

	fmt.Print("CS.reSetHeartbeatTimer\n")
}

func (c *catheState) AppendLastMsg() error {
	//ctx *context.Context, connId int64, downlinkMsg pb.DownlinkMsg
	return nil
}
