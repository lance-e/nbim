package timer

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type TimeWheel struct {
	sync.Once                             //单例工具
	interval     time.Duration            //间隔
	timer        *time.Ticker             //时间轮定时器
	slots        []*list.List             //所有槽位,每个槽位是一个双向链表,存放多个定时任务(同一个slot,但是不一定是同一个cycle)
	curSlot      int                      //当前槽位的索引
	keyToTask    map[string]*list.Element //对定时任务的映射
	stopCh       chan struct{}            //停止时间轮
	addTaskCh    chan *TaskElement        //添加定时任务
	removeTaskCh chan string              //删除定时任务
}

type TaskElement struct {
	task  func() //定时人物逻辑的闭包
	pos   int    //在时间槽中的索引
	cycle int    //轮次
	key   string //任务唯一标识
}

func NewTimeWheel(slot int, interval time.Duration) *TimeWheel {
	if slot <= 0 {
		slot = 10
	}
	if interval <= 0 {
		interval = time.Second
	}
	tw := &TimeWheel{
		interval:     interval,
		timer:        time.NewTicker(interval),
		slots:        make([]*list.List, 0, slot),
		curSlot:      0,
		keyToTask:    make(map[string]*list.Element),
		stopCh:       make(chan struct{}),
		addTaskCh:    make(chan *TaskElement),
		removeTaskCh: make(chan string),
	}

	for i := 0; i < slot; i++ {
		tw.slots = append(tw.slots, list.New())
	}

	//异步启动时间轮的常驻协程
	go tw.run()

	return tw
}

func (tw *TimeWheel) run() {
	//兜底
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("run panic happened\n")
		}
	}()

	for {
		select {
		case <-tw.stopCh: //停止定时器
			return
		case <-tw.timer.C: //批量处理定时任务
			tw.tick()
		case task := <-tw.addTaskCh: //添加定时任务
			tw.addTask(task)
		case key := <-tw.removeTaskCh: //删除定时任务
			tw.removeTask(key)
		}
	}
}

//------------public

func (tw *TimeWheel) Stop() {
	//单例工具，保证时间轮只被关闭一次
	tw.Do(func() {
		tw.timer.Stop()
		tw.stopCh <- struct{}{}
	})
}
func (tw *TimeWheel) AddTask(key string, task func(), excuteAt time.Time) {
	p, c := tw.getPosAndCycle(excuteAt)
	tw.addTaskCh <- &TaskElement{
		task:  task,
		key:   key,
		pos:   p,
		cycle: c,
	}
}

func (tw *TimeWheel) RemoveTask(key string) {
	tw.removeTaskCh <- key
}

//------------private

func (tw *TimeWheel) tick() {
	//执行完后移动到下一个槽位
	defer tw.cycleIncre()
	//获取到当前时间槽中的定时任务列表
	list := tw.slots[tw.curSlot]
	//执行当前cycle轮次的定时任务
	tw.excute(list)
}

func (tw *TimeWheel) cycleIncre() {
	tw.curSlot = (tw.curSlot + 1) % len(tw.slots)
}

func (tw *TimeWheel) excute(list *list.List) {
	//遍历列表
	for e := list.Front(); e != nil; {
		task, _ := e.Value.(*TaskElement)
		//不是该轮的任务
		if task.cycle > 0 {
			task.cycle--
			e = e.Next()
			continue
		}

		//启动一个协程异步执行定时任务
		go func() {
			//兜底
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("excute panic happened\n")
				}
			}()

			task.task()
		}()

		//清除该定时任务
		next := e.Next()
		list.Remove(e)
		delete(tw.keyToTask, task.key)
		e = next
	}
}

func (tw *TimeWheel) addTask(task *TaskElement) {
	list := tw.slots[task.pos]
	if _, ok := tw.keyToTask[task.key]; ok {
		tw.removeTask(task.key)
	}
	le := list.PushBack(task)
	tw.keyToTask[task.key] = le
}

func (tw *TimeWheel) removeTask(key string) {
	elment, ok := tw.keyToTask[key]
	if !ok {
		return
	}
	te, _ := elment.Value.(*TaskElement)
	tw.slots[te.pos].Remove(elment)
}
func (tw *TimeWheel) getPosAndCycle(excuteAt time.Time) (int, int) {
	delay := int(time.Until(excuteAt))
	//定时任务的轮次
	cycle := delay / (int(tw.interval) * len(tw.slots))
	//定时任务的索引
	pos := (tw.curSlot + delay/int(tw.interval)) % len(tw.slots)
	return pos, cycle
}
