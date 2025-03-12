package state

import (
	"context"
	"errors"
	"fmt"
	"nbim/configs"
	"nbim/pkg/db"
	"nbim/pkg/protocol/pb"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"
)

var CS *catheState

type catheState struct {
	ConnIdToConnState sync.Map     //connID到connState的映射
	Server            *StateServer //对rpc server的反引用
}

func InitCatheState() {
	CS = &catheState{
		ConnIdToConnState: sync.Map{},
		Server:            &StateServer{cmdchannel: make(chan *cmdContext, 2048)},
	}
}

func GetSlot(connId int64) int64 {
	return connId % int64(configs.GlobalConfig.StateCatheSlotNum)
}

// 登出，清除连接状态
func (c *catheState) connLogout(ctx *context.Context, connID int64) {
	if state, ok := c.ConnIdToConnState.Load(connID); ok {
		connS := state.(*connState)
		//原先的连接关闭的两种情况：
		//1.客户端主动关闭:该连接的心跳定时器未关
		//2.服务端主动关闭: 心跳超时，定时器超时才关闭的连接
		Wheel.RemoveTask(connS.heartbeatTask.Key)
		//创建重连定时器
		Wheel.AddTask(connS.reconnectTask, time.Now().Add(10*time.Second))
	}
}

// 登陆:只有登陆了才会创建连接相关状态
func (c *catheState) connLogin(ctx *context.Context, deviceId int64, connID int64) error {
	//TODO:路由信息等其他状态信息
	//缓存状态初始化
	db.RedisCli.Set(fmt.Sprintf(db.DeviceIdToConnId, deviceId), connID, db.TTL7Day)
	//本地状态保存
	c.ConnIdToConnState.Store(connID, NewConnState(connID, deviceId))
	return nil
}

func (c *catheState) connReconn(ctx *context.Context, oldConnID, newConnID int64) error {
	if state, ok := c.ConnIdToConnState.LoadAndDelete(oldConnID); ok {
		fmt.Print("connReconn: 资源还在准备复用\n")
		connS := state.(*connState)
		Wheel.RemoveTask(connS.reconnectTask.Key) //删除重连定时任务
		//复用连接状态
		connS.connID = newConnID                       //更换为新分配的connID
		c.ConnIdToConnState.Store(connS.connID, connS) //添加到映射中

		//更新redis缓存
		db.RedisCli.Set(fmt.Sprintf(db.DeviceIdToConnId, connS.deviceID), connS.connID, db.TTL7Day)

		//启动心跳定时器
		Wheel.AddTask(connS.heartbeatTask, time.Now().Add(5*time.Second))
		fmt.Printf("重连成功!\n")
		return nil
	} else {
		fmt.Printf("重连失败!\n")
		return errors.New("connection state already clear\n")
	}
}

func (c *catheState) connAck(ctx *context.Context, connID int64, sessionID uint64, msgID int64) {
	if state, ok := c.ConnIdToConnState.Load(connID); ok {
		connS := state.(*connState)

		//用msgLock来确保ack最后一条消息
		if connS.messageLock != fmt.Sprintf("%d_%d", sessionID, msgID) {
			return
		}

		//删除缓存
		key := fmt.Sprintf(db.LastMessageKey, GetSlot(connID), connID)
		if cmd := db.RedisCli.Del(key); cmd.Err() != nil {
			return
		}

		//关闭等待ack的定时任务
		if connS.messageTask != nil {
			Wheel.RemoveTask(connS.messageTask.Key)
		}
	}
}

func (c *catheState) compareAndIncrementClientID(ctx *context.Context, connID int64, oldClientID int64, sessionID uint64) bool {
	slot := GetSlot(connID)
	key := fmt.Sprintf(db.MaxClientIDKey, slot, connID, sessionID)
	//TODO:lua脚本优化,实现原子性
	//FIXME:未保证原子性，出错就不恢复
	result, err := db.RedisCli.Exists(key).Result()
	if err != nil {
		return false
	}
	//不存在
	if result == 0 {
		db.RedisCli.Set(key, 0, db.TTL7Day)
	}
	//获取value
	value, err := db.RedisCli.Get(key).Int()
	if err != nil {
		return false
	}
	//对比
	if value == int(oldClientID) {
		//自增
		err = db.RedisCli.Incr(key).Err()
		if err != nil {
			return false
		}
		//过期时间
		err = db.RedisCli.Expire(key, db.TTL7Day).Err()
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}

// 重置心跳计时器
func (c *catheState) reSetHeartbeatTimer(connId int64) {
	if state, ok := c.ConnIdToConnState.Load(connId); ok {
		connS := state.(*connState)
		Wheel.RemoveTask(connS.heartbeatTask.Key)
		Wheel.AddTask(connS.heartbeatTask, time.Now().Add(5*time.Second))
	}
}

func (c *catheState) AppendLastMsg(ctx context.Context, connId int64, downlinkMsg *pb.DownlinkMsg) error {
	state, ok := c.ConnIdToConnState.Load(connId)
	if !ok {
		return errors.New("connection state is nil")
	}
	key := fmt.Sprintf(db.LastMessageKey, GetSlot(connId), connId)
	downData, _ := proto.Marshal(downlinkMsg)

	connS := state.(*connState)
	//设置msgLock,确保是目标最新的消息
	connS.messageLock = fmt.Sprintf("%d_%d", downlinkMsg.SessionId, downlinkMsg.Seq)
	//设置消息定时器，直到被收到下行消息的ack
	connS.SetMessageTimer()

	//存入redis
	if cmd := db.RedisCli.Set(key, downData, db.TTL7Day); cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (c *catheState) GetLastMsg(ctx context.Context, connId int64) (*pb.DownlinkMsg, error) {
	key := fmt.Sprintf(db.LastMessageKey, GetSlot(connId), connId)
	cmd := db.RedisCli.Get(key)
	if cmd == nil {
		return nil, errors.New("rediscmd is nil")
	}
	data, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	down := &pb.DownlinkMsg{}
	err = proto.Unmarshal(data, down)

	return down, err
}
