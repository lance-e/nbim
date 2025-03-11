package id_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ... (雪花算法代码)

func TestSnowflake(t *testing.T) {
	// 正常情况
	snowflake, err := NewSnowflake(1)
	assert.NoError(t, err)

	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := snowflake.Generate()
		assert.NotZero(t, id)
		_, ok := ids[id]
		assert.False(t, ok, "duplicate ID generated")
		ids[id] = true
	}

	// 并发生成
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snowflake, err := NewSnowflake(1)
			assert.NoError(t, err)
			for j := 0; j < 100; j++ {
				id := snowflake.Generate()
				assert.NotZero(t, id)
			}
		}()
	}
	wg.Wait()

	// workerID 范围
	_, err = NewSnowflake(-1)
	assert.Error(t, err, "worker ID out of range")

	_, err = NewSnowflake(maxNode + 1)
	assert.Error(t, err, "worker ID out of range")

	snowflake, err = NewSnowflake(maxNode)
	assert.NoError(t, err)
	id := snowflake.Generate()
	assert.NotZero(t, id)

	// 序列号溢出
	snowflake, err = NewSnowflake(1)
	assert.NoError(t, err)

	for i := 0; i <= maxSequence; i++ {
		snowflake.Generate()
	}

	// 时钟回拨
	snowflake.lastTimestamp = time.Now().UnixNano()/1e6 + 1000 // 模拟时钟回拨
	assert.Panics(t, func() {
		snowflake.Generate()
	}, "clock is moving backwards")
}

const (
	nodeBit        = 10
	sequenceBit    = 12
	maxNode        = -1 ^ (-1 << nodeBit)     //最大节点数
	maxSequence    = -1 ^ (-1 << sequenceBit) //最大序列号数
	nodeShift      = 12                       //节点左移12位
	timestampShift = 22                       //时间戳左移22位
)

type Snowflake struct {
	mutex         sync.Mutex
	lastTimestamp int64 //上一个毫秒时间戳
	nodeId        int64 //工作机器id , 最大 2 ^ 10
	sequence      int64 //记录已经分配的sequence , 最大2 ^ 12
}

func NewSnowflake(nodeId int64) (*Snowflake, error) {
	if nodeId < 0 || nodeId > maxNode {
		return nil, errors.New("nodeId is wrong")
	}
	return &Snowflake{
		mutex:         sync.Mutex{},
		lastTimestamp: time.Now().UnixNano() / 1e6, //毫秒时间戳
		nodeId:        nodeId,
		sequence:      0,
	}, nil
}

// Generate:生成唯一id
func (s *Snowflake) Generate() int64 {
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
