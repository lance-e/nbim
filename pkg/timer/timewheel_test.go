package timer_test

import (
	"nbim/pkg/timer"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeWheel(t *testing.T) {
	// 基本功能测试
	tw := timer.NewTimeWheel(10, time.Second)
	defer tw.Stop()

	var executed bool
	tw.AddTask("test1", func() {
		executed = true
	}, time.Now().Add(2*time.Second))

	time.Sleep(3 * time.Second)
	assert.True(t, executed, "task should be executed")

	// 删除任务测试
	executed = false
	tw.AddTask("test2", func() {
		executed = true
	}, time.Now().Add(2*time.Second))
	tw.RemoveTask("test2")
	time.Sleep(3 * time.Second)
	assert.False(t, executed, "task should not be executed")

	// 停止测试
	tw2 := timer.NewTimeWheel(10, time.Second)
	tw2.Stop()
	time.Sleep(time.Second) // 确保停止协程执行完毕

	// 边界情况测试
	tw3 := timer.NewTimeWheel(0, 0)
	assert.NotNil(t, tw3, "should create time wheel with default values")
	tw3.Stop()

	tw4 := timer.NewTimeWheel(10, time.Second)
	var executed2 bool
	tw4.AddTask("test3", func() {
		executed2 = true
	}, time.Now().Add(2*time.Second))
	tw4.AddTask("test3", func() {}, time.Now().Add(3*time.Second)) //覆盖之前的任务
	time.Sleep(3 * time.Second)
	assert.False(t, executed2, "old task should not be executed")
	tw4.Stop()

	tw5 := timer.NewTimeWheel(10, time.Second)
	tw5.RemoveTask("notexist")
	tw5.Stop()

	// 并发测试
	tw6 := timer.NewTimeWheel(10, time.Millisecond*100)
	defer tw6.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tw6.AddTask(strconv.Itoa(id), func() {}, time.Now().Add(time.Millisecond*500))
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)

	//轮次测试
	tw7 := timer.NewTimeWheel(5, time.Millisecond*100)
	defer tw7.Stop()

	var executed3 bool
	tw7.AddTask("test4", func() {
		executed3 = true
	}, time.Now().Add(time.Millisecond*600)) //需要跨一轮

	time.Sleep(time.Second)
	assert.True(t, executed3, "task should be executed after one cycle")

	// 异常测试
	tw8 := timer.NewTimeWheel(10, time.Second)
	defer tw8.Stop()

	tw8.AddTask("test5", func() {
		panic("test panic")
	}, time.Now().Add(1*time.Second))

	time.Sleep(2 * time.Second) // 确保有足够时间让panic发生
}
