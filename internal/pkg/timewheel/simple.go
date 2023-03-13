package timewheel

import (
	"log"
	"sync"
	"time"

	"github.com/sourcegraph/conc/pool"
)

// SimpleTimeWheel 简单时间轮
type SimpleTimeWheel struct {
	interval  time.Duration
	ticker    *time.Ticker
	slot      []*slot
	indicator *sync.Map
	tickIndex int
	onTick    SimpleHandler
	taskChan  chan any
	quitChan  chan any
}

// SimpleHandler 处理函数
type SimpleHandler func(*SimpleTimeWheel, string, any)

func NewSimpleTimeWheel(delay time.Duration, numSlot int, handler SimpleHandler) *SimpleTimeWheel {
	timeWheel := &SimpleTimeWheel{
		taskChan:  make(chan any, 100),
		quitChan:  make(chan any),
		indicator: &sync.Map{},
		interval:  delay,
		ticker:    time.NewTicker(delay),
		onTick:    handler,
	}

	for i := 0; i < numSlot; i++ {
		timeWheel.slot = append(timeWheel.slot, newSlot(i))
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
		case v := <-t.taskChan:
			el, ok := v.(*element)
			if !ok {
				continue
			}

			circleSlot := t.slot[t.getCircleAndSlot(el)]
			circleSlot.add(el)
			t.indicator.Store(el.key, circleSlot)
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
			slot.elements.Range(func(key, value any) bool {
				el, ok := value.(*element)
				if !ok {
					return true
				}

				t.indicator.Delete(el.key)
				slot.remove(el.key)

				worker.Go(func() {
					if el.expire <= time.Now().Unix() {
						t.onTick(t, el.key, el.value)
					} else {
						second := el.expire - time.Now().Unix()
						if err := t.Add(el.key, el.value, time.Duration(second)*time.Second); err != nil {
							log.Printf("时间轮降级失败 err:%s", err.Error())
						}
					}
				})

				return true
			})
		}
	}
}

// Add 添加任务
func (t *SimpleTimeWheel) Add(key string, task any, delay time.Duration) error {

	t.taskChan <- &element{key: key, value: task, expire: time.Now().Add(delay).Unix()}

	return nil
}

func (t *SimpleTimeWheel) Remove(key string) {
	if value, ok := t.indicator.Load(key); ok {
		if slot, ok := value.(*slot); ok {
			slot.remove(key)
			t.indicator.Delete(key)
		}
	}
}

func (t *SimpleTimeWheel) getCircleAndSlot(el *element) int {

	remainingTime := int(el.expire - time.Now().Unix())
	if remainingTime <= 0 {
		remainingTime = 0
	}

	return (t.tickIndex + remainingTime/int(t.interval.Seconds())) % len(t.slot)
}
