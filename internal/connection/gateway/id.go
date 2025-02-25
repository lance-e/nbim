package gateway

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

type snowflake struct {
	mutex         sync.Mutex
	lastTimestamp int64 //上一个毫秒时间戳
	nodeId        int64 //工作机器id , 最大 2 ^ 10
	sequence      int64 //记录已经分配的sequence , 最大2 ^ 12
}

var Snowflake *snowflake

func NewSnowflake(nodeId int64) (*snowflake, error) {
	if nodeId < 0 || nodeId > maxNode {
		return nil, errors.New("nodeId is wrong")
	}
	return &snowflake{
		mutex:         sync.Mutex{},
		lastTimestamp: time.Now().UnixNano() / 1e6, //毫秒时间戳
		nodeId:        nodeId,
		sequence:      0,
	}, nil
}

// Generate:生成唯一id
func (s *snowflake) Generate() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now < s.lastTimestamp { //时间回拨，
		panic("close is moving backwards")
	}
	if now == s.lastTimestamp { //此次时间戳与上一次相同
		s.sequence = (s.sequence + 1) & maxSequence //序列号+1 ，同时控制在范围内
		if s.sequence == 0 {                        //	说明溢出了
			for now <= s.lastTimestamp {
				//等待到下一个毫秒再继续分配，虽然序列号可能和前面的重复，但是时间戳不会，保证了序列号全局唯一
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else { //新的时间戳
		s.sequence = 0
	}
	s.lastTimestamp = now //更新时间戳
	id := now<<timestampShift | s.nodeId<<nodeShift | s.sequence
	return id
}
