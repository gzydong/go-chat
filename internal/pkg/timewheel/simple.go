package timewheel

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sourcegraph/conc/pool"
)

// SimpleTimeWheel 简单时间轮
type SimpleTimeWheel struct {
	interval  time.Duration
	ticker    *time.Ticker
	slot      []cmap.ConcurrentMap[string, *element]
	indicator cmap.ConcurrentMap[string, int]
	tickIndex int
	onTick    SimpleHandler
	taskChan  chan *element
	quitChan  chan struct{}
}

// SimpleHandler 处理函数
type SimpleHandler func(*SimpleTimeWheel, string, any)

func NewSimpleTimeWheel(delay time.Duration, numSlot int, handler SimpleHandler) *SimpleTimeWheel {
	timeWheel := &SimpleTimeWheel{
		taskChan:  make(chan *element, 100),
		quitChan:  make(chan struct{}),
		indicator: cmap.New[int](),
		interval:  delay,
		ticker:    time.NewTicker(delay),
		onTick:    handler,
	}

	for i := 0; i < numSlot; i++ {
		timeWheel.slot = append(timeWheel.slot, cmap.New[*element]())
	}

	return timeWheel
}

// Start 启动时间轮任务
func (t *SimpleTimeWheel) Start() {

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

func (t *SimpleTimeWheel) Stop() {
	close(t.quitChan)
}

func (t *SimpleTimeWheel) run() {

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
				el := item.Val

				slot.Remove(el.key)
				t.indicator.Remove(el.key)

				worker.Go(func() {
					if el.expire <= time.Now().Unix() {
						t.onTick(t, el.key, el.value)
					} else {
						second := el.expire - time.Now().Unix()
						_ = t.Add(el.key, el.value, time.Duration(second)*time.Second)
					}
				})
			}
		}
	}
}

// Add 添加任务
func (t *SimpleTimeWheel) Add(key string, task any, delay time.Duration) error {

	t.taskChan <- &element{key: key, value: task, expire: time.Now().Add(delay).Unix()}

	return nil
}

func (t *SimpleTimeWheel) Remove(key string) {
	if value, ok := t.indicator.Get(key); ok {
		t.slot[value].Remove(key)
		t.indicator.Remove(key)
	}
}

func (t *SimpleTimeWheel) getCircleAndSlot(el *element) int {

	remainingTime := int(el.expire - time.Now().Unix())
	if remainingTime <= 0 {
		remainingTime = 0
	}

	return (t.tickIndex + remainingTime/int(t.interval.Seconds())) % len(t.slot)
}
