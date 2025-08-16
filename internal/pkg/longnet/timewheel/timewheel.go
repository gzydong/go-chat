package timewheel

import (
	"container/list"
	"errors"
	"hash/fnv"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var taskPool = sync.Pool{
	New: func() interface{} {
		return &Task{}
	},
}

// Task 结构体
type Task struct {
	ID       int64
	Callback func(id int64)
}

func NewTask(id int64, cb func(id int64)) *Task {
	obj := taskPool.Get()
	var pooledTask *Task

	if t, ok := obj.(*Task); ok && t != nil {
		pooledTask = t
		pooledTask.Reset()
	} else {
		pooledTask = &Task{}
	}

	pooledTask.ID = id
	pooledTask.Callback = cb
	return pooledTask
}

func (t *Task) Reset() {
	t.ID = 0
	t.Callback = nil
}

// Bucket 是线程安全的任务容器
type Bucket struct {
	l  *list.List
	mu sync.Mutex
}

func NewBucket() *Bucket {
	return &Bucket{
		l: list.New(),
	}
}

func (b *Bucket) Add(task *Task) *list.Element {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.l.PushBack(task)
}

func (b *Bucket) Remove(ele *list.Element) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.l.Remove(ele)
}

func (b *Bucket) Execute(cb func(taskId int64), consumer chan *Task) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for e := b.l.Front(); e != nil; {
		next := e.Next()
		task := e.Value.(*Task)

		cb(task.ID)
		consumer <- task

		b.l.Remove(e)
		e = next
	}
}

func (b *Bucket) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.l.Init()
}

type TaskElement struct {
	slotIdx int64
	element *list.Element
}

// TimingWheel 时间轮核心结构
type TimingWheel struct {
	interval    time.Duration
	numSlots    int
	ticker      *time.Ticker
	slots       []*Bucket
	currentPos  int64
	stopChan    chan struct{}
	mu          sync.RWMutex
	consumer    chan *Task
	callback    func(taskId int64)
	closed      atomic.Bool
	globalTasks cmap.ConcurrentMap[int64, *TaskElement]
}

func NewTimingWheel(interval time.Duration, numSlots int, worker int) *TimingWheel {
	if interval <= 0 || numSlots <= 0 {
		panic("interval and numSlots must be positive")
	}

	if worker <= 0 {
		panic("worker must be positive")
	}

	tw := &TimingWheel{
		interval:    interval,
		numSlots:    numSlots,
		slots:       make([]*Bucket, numSlots),
		currentPos:  0,
		stopChan:    make(chan struct{}),
		globalTasks: cmap.NewWithCustomShardingFunction[int64, *TaskElement](fnv32),
		consumer:    make(chan *Task, 1000),
	}

	for i := 0; i < numSlots; i++ {
		tw.slots[i] = NewBucket()
	}

	for i := 0; i < worker; i++ {
		go tw.startConsumer()
	}

	tw.ticker = time.NewTicker(interval)
	go tw.run()

	return tw
}

func (tw *TimingWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.advanceAndRun()
		case <-tw.stopChan:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimingWheel) advanceAndRun() {
	var idx int64
	var slot *Bucket

	// 获取当前槽索引
	tw.mu.RLock()
	idx = tw.currentPos % int64(tw.numSlots)

	slot = tw.slots[idx]
	tw.mu.RUnlock()

	go slot.Execute(tw.removeGlobalTask, tw.consumer)

	// 更新指针位置
	tw.mu.Lock()
	tw.currentPos++
	tw.mu.Unlock()
}

func (tw *TimingWheel) SetCallback(cb func(taskId int64)) {
	tw.callback = cb
}

// AddTask 添加一个延迟执行的任务
func (tw *TimingWheel) AddTask(taskId int64, d time.Duration, callback func(id int64)) error {
	if d < 0 {
		return errors.New("delay must be positive")
	}

	maxDelay := time.Duration(tw.numSlots) * tw.interval
	if d >= maxDelay {
		return errors.New("delay exceeds wheel capacity")
	}

	delayTicks := int64(d / tw.interval)
	slotIdx := (tw.currentPos + delayTicks) % int64(tw.numSlots)

	task := NewTask(taskId, callback)

	// 如果已存在相同 ID 的任务，先移除
	info, exists := tw.globalTasks.Get(taskId)
	if exists {
		tw.slots[info.slotIdx].Remove(info.element)
	}

	// 添加到目标槽
	ele := tw.slots[slotIdx].Add(task)

	// 记录全局信息
	tw.globalTasks.Set(taskId, &TaskElement{slotIdx: slotIdx, element: ele})
	return nil
}

// Cancel 根据 ID 取消任意槽中的任务
func (tw *TimingWheel) Cancel(taskId int64) {
	task, exists := tw.globalTasks.Get(taskId)
	if !exists {
		return
	}

	tw.slots[task.slotIdx].Remove(task.element)
	tw.globalTasks.Remove(taskId)
}

func (tw *TimingWheel) removeGlobalTask(taskId int64) {
	tw.globalTasks.Remove(taskId)
}

// Stop 停止时间轮
func (tw *TimingWheel) Stop() {
	if tw.closed.Swap(true) {
		close(tw.stopChan)
		close(tw.consumer)
	}
}

func (tw *TimingWheel) startConsumer() {
	for task := range tw.consumer {
		if tw.closed.Load() {
			return
		}

		if task.Callback != nil {
			safeCallback(func() { task.Callback(task.ID) })
		} else if tw.callback != nil {
			safeCallback(func() { tw.callback(task.ID) })
		}

		task.Callback = nil
		task.ID = 0
		taskPool.Put(task)
	}
}

func fnv32(key int64) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(strconv.FormatInt(key, 10)))
	return h.Sum32()
}

func safeCallback(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in task callback: %v\n", r)
		}
	}()

	fn()
}
