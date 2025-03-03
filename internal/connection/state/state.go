package state

import (
	"context"
	"fmt"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"nbim/pkg/timer"
	"time"
)

// 全局时间轮
var Wheel timer.TimeWheel

// 连接状态
type connState struct {
	connID        int64
	deviceID      int64
	messageTask   *timer.TaskElement //消息定时任务
	reconnectTask *timer.TaskElement //重连定时任务
	heartbeatTask *timer.TaskElement //心跳定时任务
}

func NewConnState(connId int64, deviceId int64) *connState {
	state := &connState{
		connID:   connId,
		deviceID: deviceId,
	}

	//-------------reconnect
	state.reconnectTask = timer.NewTaskElement(fmt.Sprintf("%d|reconnect", connId), func() {
		//重连超时，清除连接状态
		fmt.Print("滴滴滴，重连定时器超时\n")
		state.clear()
	})
	//-------------

	//-------------heartbeat
	state.heartbeatTask = timer.NewTaskElement(fmt.Sprintf("%d|heartbeat", connId), func() {
		fmt.Print("滴滴滴，心跳定时器超时\n")
		//心跳超时,强制断开连接
		rpc.GetGatewayClient().CloseConn(context.TODO(), &pb.GatewayRequest{
			ConnId: connId,
		})
		//主动调用closecConn rpc直接就断开了连接，不会再经过onconnnetion回调
		//需要启动重连定时器和执行清除状态逻辑
		Wheel.AddTask(state.reconnectTask, time.Now().Add(10*time.Second))
	})
	//-------------

	//全局时间轮中添加心跳任务
	Wheel.AddTask(state.heartbeatTask, time.Now().Add(5*time.Second))
	return state
}

func (c *connState) SetMessageTimer() {
	//清除之前的定时任务
	if c.messageTask != nil {
		Wheel.RemoveTask(c.messageTask.Key)
		c.messageTask = nil
	}
	c.messageTask = timer.NewTaskElement(fmt.Sprintf("%d|message", c.connID), func() {
		fmt.Printf("滴滴滴，我的ack怎么还没到?\n")
		//重发消息
		reSendDownlinkMessage(c.connID)
	})

	Wheel.AddTask(c.messageTask, time.Now().Add(100*time.Millisecond))
}

// clearConnState:执行到该函数说明之前的超时任务已经执行，直接清除
func (c *connState) clear() {
	if c.messageTask != nil {
		Wheel.RemoveTask(c.messageTask.Key)
	}
	if c.reconnectTask != nil {
		Wheel.RemoveTask(c.reconnectTask.Key)
	}
	if c.heartbeatTask != nil {
		Wheel.RemoveTask(c.heartbeatTask.Key)
	}
	//TODO:清除其他状态,memory ,redis,router等
	CS.ConnIdToConnState.Delete(c.connID)
	fmt.Printf("已清除[%d]所有状态!\n",c.connID)               
}
