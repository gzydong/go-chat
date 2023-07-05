package timewheel

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sourcegraph/conc/pool"
)

type element struct {
	key    string // 任务key
	value  any    // 任务信息
	expire int64  // 过期时间
}

type slot struct {
	id       int       // 轮次ID
	elements *sync.Map // 成员列表
}

func newSlot(id int) *slot {
	s := &slot{id: id}
	s.elements = &sync.Map{}
	return s
}

func (s *slot) add(el *element) {
	s.elements.LoadOrStore(el.key, el)
}

func (s *slot) remove(key any) {
	s.elements.Delete(key)
}

// 时间轮环
type circle struct {
	index     int          // 轮环ID
	tickIndex int          // 当前插槽
	ticker    *time.Ticker // 计时器
	slot      []*slot      // 插槽列表
}

func newCircle(index int, numSlots int, ticker *time.Ticker) *circle {
	c := &circle{
		index:     index,
		tickIndex: 0,
		ticker:    ticker,
		slot:      make([]*slot, 0, numSlots),
	}

	for i := 0; i < numSlots; i++ {
		c.slot = append(c.slot, newSlot(i))
	}

	return c
}

// Handler 处理函数
type Handler func(*TimeWheel, any)

// TimeWheel 分层时间轮
// 第一层秒  0 ~ 59
// 第二层分  0 ~ 59
// 第三场时  0 ~ 23
// @see https://blog.csdn.net/daocaokafei/article/details/126456817
type TimeWheel struct {
	circle    []*circle
	onTick    Handler
	taskChan  chan any
	quitChan  chan any
	indicator *sync.Map
}

func NewTimeWheel(handler Handler) *TimeWheel {

	timeWheel := &TimeWheel{
		taskChan:  make(chan any, 100),
		quitChan:  make(chan any),
		indicator: &sync.Map{},
		onTick:    handler,
	}

	// 初始化时间轮
	timeWheel.circle = []*circle{
		newCircle(0, 60, time.NewTicker(time.Second)),
		newCircle(1, 60, time.NewTicker(time.Minute)),
		newCircle(2, 24, time.NewTicker(time.Hour)),
	}

	return timeWheel
}

// Start 启动时间轮任务
func (t *TimeWheel) Start() {
	defer fmt.Println("TimeWheel Stop")

	// 协程启动3个时间分层轮
	for _, c := range t.circle {
		go func(c *circle) {
			t.runTimeWheel(c)
		}(c)
	}

	for {
		select {
		case <-t.quitChan:
			return
		case v := <-t.taskChan:
			el, ok := v.(*element)
			if !ok {
				continue
			}

			circleIndex, slotIndex := t.getCircleAndSlot(el)

			circleSlot := t.circle[circleIndex].slot[slotIndex]
			circleSlot.add(el)
			t.indicator.Store(el.value, circleSlot)
		}
	}
}

func (t *TimeWheel) getCircleAndSlot(el *element) (int, int64) {

	var (
		circleIndex   int
		slotIndex     int64
		remainingTime = int(el.expire - time.Now().Unix())
	)

	if remainingTime <= 0 {
		remainingTime = 0
	}

	if remainingTime < 60 {
		circleIndex = 0
		slotIndex = int64((t.getCurrentTickIndex(0) + remainingTime) % 60)
	} else if int(remainingTime/60) < 60 {
		circleIndex = 1
		slotIndex = int64((t.getCurrentTickIndex(1) + remainingTime/60) % 60)
	} else {
		circleIndex = 2
		slotIndex = int64((t.getCurrentTickIndex(1) + remainingTime/3600) % 60)
	}

	// fmt.Printf("加入任务信息 过期时间:%s 当前剩余时间:%ds 加入时间轮:%d 槽位:%d \n",
	// 	time.Now().Add(time.Duration(remainingTime)*time.Second).Format("2006-01-02 15:04:05"),
	// 	remainingTime,
	// 	circleIndex,
	// 	slotIndex,
	// )

	return circleIndex, slotIndex
}

func (t *TimeWheel) runTimeWheel(circle *circle) {

	defer fmt.Printf("[%d]RunTimeWheel Stop\n", circle.index)

	worker := pool.New().WithMaxGoroutines(10)

	for {
		select {
		case <-t.quitChan:
			circle.ticker.Stop()
			return
		case <-circle.ticker.C:

			tickIndex := circle.tickIndex

			circle.tickIndex++
			if circle.tickIndex >= len(circle.slot) {
				circle.tickIndex = 0
			}

			circleSlot := circle.slot[tickIndex]

			circleSlot.elements.Range(func(_, value any) bool {
				if el, ok := value.(*element); ok {
					t.indicator.Delete(el.value)
					circleSlot.remove(el.value)

					worker.Go(func() {
						if el.expire <= time.Now().Unix() {
							t.onTick(t, el.value)
						} else {
							second := el.expire - time.Now().Unix()
							if err := t.Add(el.value, time.Duration(second)*time.Second); err != nil {
								log.Printf("时间轮降级失败 err:%s", err.Error())
							}
						}
					})
				}

				return true
			})
		}
	}
}

func (t *TimeWheel) Stop() {
	close(t.quitChan)
}

// Add 添加任务
// 注: 不支持大于24小时的延时任务
func (t *TimeWheel) Add(task any, delay time.Duration) error {

	if delay > 24*time.Hour {
		return errors.New("max delay 24 hour")
	}

	t.taskChan <- &element{value: task, expire: time.Now().Add(delay).Unix()}

	return nil
}

func (t *TimeWheel) Remove(task any) {
	if value, ok := t.indicator.Load(task); ok {
		if slot, ok := value.(*slot); ok {
			slot.remove(task)
			t.indicator.Delete(task)
		}
	}
}

func (t *TimeWheel) getCurrentTickIndex(circleIndex int) int {
	return t.circle[circleIndex].tickIndex
}
