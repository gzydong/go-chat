package timewheel

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sourcegraph/conc/pool"
)

type entry[T any] struct {
	key    string
	value  T
	expire int64
}

// SimpleTimeWheel 简单时间轮
type SimpleTimeWheel[T any] struct {
	interval  time.Duration
	ticker    *time.Ticker
	tickIndex int
	slot      []cmap.ConcurrentMap[string, *entry[T]]
	indicator cmap.ConcurrentMap[string, int]
	onTick    SimpleHandler[T]
	taskChan  chan *entry[T]
	quitChan  chan struct{}
}

// SimpleHandler 处理函数
type SimpleHandler[T any] func(*SimpleTimeWheel[T], string, T)

func NewSimpleTimeWheel[T any](delay time.Duration, numSlot int, handler SimpleHandler[T]) *SimpleTimeWheel[T] {
	timeWheel := &SimpleTimeWheel[T]{
		taskChan:  make(chan *entry[T], 100),
		quitChan:  make(chan struct{}),
		indicator: cmap.New[int](),
		interval:  delay,
		ticker:    time.NewTicker(delay),
		onTick:    handler,
	}

	for i := 0; i < numSlot; i++ {
		timeWheel.slot = append(timeWheel.slot, cmap.New[*entry[T]]())
	}

	return timeWheel
}

// Start 启动时间轮任务
func (t *SimpleTimeWheel[T]) Start() {

	go t.run()

	for {
		select {
		case <-t.quitChan:
			return
		case el := <-t.taskChan:
			t.Remove(el.key)

			slotIndex := t.getCircleAndSlot(el)
			t.slot[slotIndex].Set(el.key, el)
			t.indicator.Set(el.key, slotIndex)
		}
	}
}

func (t *SimpleTimeWheel[T]) Stop() {
	close(t.quitChan)
}

func (t *SimpleTimeWheel[T]) run() {

	worker := pool.New().WithMaxGoroutines(10)

	for {
		select {
		case <-t.quitChan:
			t.ticker.Stop()
			return
		case <-t.ticker.C:
			tickIndex := t.tickIndex

			t.tickIndex++
			if t.tickIndex >= len(t.slot) {
				t.tickIndex = 0
			}

			slot := t.slot[tickIndex]
			for item := range slot.IterBuffered() {
				v := item.Val

				slot.Remove(v.key)
				t.indicator.Remove(v.key)

				worker.Go(func() {
					unix := time.Now().Unix()
					if v.expire <= unix {
						t.onTick(t, v.key, v.value)
					} else {
						t.Add(v.key, v.value, time.Duration(v.expire-unix)*time.Second)
					}
				})
			}
		}
	}
}

// Add 添加任务
func (t *SimpleTimeWheel[T]) Add(key string, value T, delay time.Duration) {
	t.taskChan <- &entry[T]{key: key, value: value, expire: time.Now().Add(delay).Unix()}
}

func (t *SimpleTimeWheel[T]) Remove(key string) {
	if value, ok := t.indicator.Get(key); ok {
		t.slot[value].Remove(key)
		t.indicator.Remove(key)
	}
}

func (t *SimpleTimeWheel[T]) getCircleAndSlot(e *entry[T]) int {

	remainingTime := int(e.expire - time.Now().Unix())
	if remainingTime <= 0 {
		remainingTime = 0
	}

	return (t.tickIndex + remainingTime/int(t.interval.Seconds())) % len(t.slot)
}
