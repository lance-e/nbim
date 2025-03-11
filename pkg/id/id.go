package id

import (
	"errors"
	"sync"
	"time"
)

//--------------- connection ID generater -----------------
//  snowflake
/* +--------------------------------------------------------------------------+ */
/* | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID | */
/* +--------------------------------------------------------------------------+ */

const (
	nodeBit        = 10
	sequenceBit    = 12
	maxNode        = -1 ^ (-1 << nodeBit)     //最大节点数
	maxSequence    = -1 ^ (-1 << sequenceBit) //最大序列号数
	nodeShift      = 12                       //节点左移12位
	timestampShift = 22                       //时间戳左移22位
)

type Snowflake struct {
	Mutex         sync.Mutex
	LastTimestamp int64 //上一个毫秒时间戳
	NodeId        int64 //工作机器id , 最大 2 ^ 10
	Sequence      int64 //记录已经分配的sequence , 最大2 ^ 12
}

func NewSnowflake(nodeId int64) (*Snowflake, error) {
	if nodeId < 0 || nodeId > maxNode {
		return nil, errors.New("nodeId is wrong")
	}
	return &Snowflake{
		Mutex:         sync.Mutex{},
		LastTimestamp: time.Now().UnixNano() / 1e6, //毫秒时间戳
		NodeId:        nodeId,
		Sequence:      0,
	}, nil
}

// Generate:生成唯一id
func (s *Snowflake) Generate() int64 {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now < s.LastTimestamp { //时间回拨，
		panic("clock is moving backwards")
	}
	if now == s.LastTimestamp { //此次时间戳与上一次相同
		s.Sequence = (s.Sequence + 1) & maxSequence //序列号+1 ，同时控制在范围内
		if s.Sequence == 0 {                        //	说明溢出了
			for now <= s.LastTimestamp {
				//等待到下一个毫秒再继续分配，虽然序列号可能和前面的重复，但是时间戳不会，保证了序列号全局唯一
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else { //新的时间戳
		s.Sequence = 0
	}
	s.LastTimestamp = now //更新时间戳
	id := now<<timestampShift | s.NodeId<<nodeShift | s.Sequence
	return id
}
