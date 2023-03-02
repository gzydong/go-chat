package timewheel

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	WheelSecond = 60
	WheelMinute = 60
	WheelHour   = 24
)

type element struct {
	task   any   // 任务信息
	expire int64 // 过期时间
}

type slot struct {
	id       int              // 轮次ID
	elements map[any]*element // 成员列表
}

func newSlot(id int) *slot {
	s := &slot{id: id}
	s.elements = make(map[any]*element)
	return s
}

func (s *slot) add(el *element) {
	s.elements[el.task] = el
}

func (s *slot) remove(task any) {
	delete(s.elements, task)
}

// 时间轮环
type wheel struct {
	wheelIndex       int // 时间轮环ID
	currentTickIndex int // 当前插槽Index
	ticker           *time.Ticker
	slot             []*slot // 插槽列表
}

func newWheel(wheelIndex int, slotNum int, ticker *time.Ticker) *wheel {
	wheel := &wheel{
		wheelIndex:       wheelIndex,
		currentTickIndex: 0,
		ticker:           ticker,
		slot:             make([]*slot, 0),
	}

	for i := 0; i < slotNum; i++ {
		wheel.slot = append(wheel.slot, newSlot(i))
	}

	return wheel
}

// Handler 处理函数
type Handler func(*TimeWheel, any)

// TimeWheel 分层时间轮
// 第一层秒  0 ~ 59
// 第一场分  0 ~ 59
// 第一场时  0 ~ 23
// @see https://blog.csdn.net/daocaokafei/article/details/126456817
type TimeWheel struct {
	wheel     []*wheel
	lock      sync.RWMutex
	indicator map[any]*slot
	onTick    Handler
	taskChan  chan any
	quitChan  chan any
}

func NewTimeWheel(handler Handler) *TimeWheel {

	timeWheel := &TimeWheel{
		taskChan:  make(chan any),
		quitChan:  make(chan any),
		indicator: make(map[any]*slot, 0),
		onTick:    handler,
	}

	// 初始化时间轮
	timeWheel.wheel = []*wheel{
		newWheel(0, 60, time.NewTicker(time.Second)), // 秒-时间环
		newWheel(1, 60, time.NewTicker(time.Minute)), // 分-时间环
		newWheel(2, 24, time.NewTicker(time.Hour)),   // 时-时间环
	}

	return timeWheel
}

// Start 启动时间轮任务
func (t *TimeWheel) Start() {
	defer fmt.Println("TimeWheel Stop")

	// 协程启动3个时间分层轮
	for i := range t.wheel {
		wheel := t.wheel[i]
		go t.runTimeWheel(wheel.wheelIndex, wheel.ticker)
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

			// 根据任务的过期时间，计算任务所在的时间轮及时间轮插槽位
			currentWheelIndex, currentTickIndex := t.getWheelIndex(el)

			// 过去插槽节点
			wheelSlot := t.wheel[currentWheelIndex].slot[currentTickIndex]

			t.lock.Lock()
			wheelSlot.add(el)
			t.indicator[el.task] = wheelSlot
			t.lock.Unlock()
		}
	}
}

func (t *TimeWheel) getWheelIndex(el *element) (int, int64) {

	var (
		currentWheelIndex int
		currentTickIndex  int64
		remainingTime     = el.expire - time.Now().Unix()
	)

	if remainingTime <= 0 {
		remainingTime = 0
	}

	if remainingTime < 60 {
		// 加入秒时间轮
		currentWheelIndex = 0
		currentTickIndex = int64(t.getCurrentTickIndex(0)) + remainingTime

		if currentTickIndex >= WheelSecond {
			currentTickIndex = currentTickIndex - WheelSecond
		}
	} else if int(remainingTime/60) < 60 {
		// 加入分时间轮
		currentWheelIndex = 1
		currentTickIndex = int64(t.getCurrentTickIndex(1)) + (remainingTime / 60) - 1

		if currentTickIndex >= WheelMinute {
			currentTickIndex = currentTickIndex - WheelMinute
		}
	} else {
		// 加入时时间轮
		currentWheelIndex = 2
		currentTickIndex = int64(t.getCurrentTickIndex(2)) + (remainingTime / (60 * 60)) - 1

		if currentTickIndex >= WheelHour {
			currentTickIndex = currentTickIndex - WheelHour
		}
	}

	fmt.Printf("加入任务信息 过期时间:%s 当前剩余时间:%ds 加入时间轮:%d 槽位:%d \n",
		time.Now().Add(time.Duration(remainingTime)*time.Second).Format("2006-01-02 15:04:05"),
		remainingTime,
		currentWheelIndex,
		currentTickIndex,
	)

	return currentWheelIndex, currentTickIndex
}

func (t *TimeWheel) runTimeWheel(wheelIndex int, ticker *time.Ticker) {

	defer fmt.Printf("[%d]RunTimeWheel Stop\n", wheelIndex)

	for {
		select {
		case <-t.quitChan:
			ticker.Stop()
			return
		case <-ticker.C:
			// 取出对应的时间轮信息
			timeWheel := t.wheel[wheelIndex]

			currentTickIndex := timeWheel.currentTickIndex

			// 累增当前index
			timeWheel.currentTickIndex++
			if timeWheel.currentTickIndex >= len(timeWheel.slot) {
				timeWheel.currentTickIndex = 0
			}

			wheelSlot := timeWheel.slot[currentTickIndex]
			for _, v := range wheelSlot.elements {

				t.lock.Lock()
				wheelSlot.remove(v.task)
				delete(t.indicator, v)
				t.lock.Unlock()

				if v.expire <= time.Now().Unix() {
					t.onTick(t, v.task)
				} else {
					second := v.expire - time.Now().Unix()
					if err := t.Add(v.task, time.Duration(second)*time.Second); err != nil {
						log.Printf("分级别时间轮降级为妙级别时间轮失败 err:%s", err.Error())
					}

					// fmt.Printf("此任务需要降级处理: %d remainingTime: %d \n", v.expire, second)
				}
			}
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

	// fmt.Println("创建时间:", time.Now().Format("2006-01-02 15:04:05"), "过期时间:", time.Now().Add(delay).Format("2006-01-02 15:04:05"))
	t.taskChan <- &element{task: task, expire: time.Now().Add(delay).Unix()}

	return nil
}

func (t *TimeWheel) Remove(task any) {
	if slot, ok := t.indicator[task]; ok {
		t.lock.Lock()
		slot.remove(task)
		delete(t.indicator, task)
		t.lock.Unlock()
	}
}

func (t *TimeWheel) getCurrentTickIndex(wheelIndex int) int {
	return t.wheel[wheelIndex].currentTickIndex
}
