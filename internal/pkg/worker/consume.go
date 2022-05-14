package worker

import (
	"sync"
)

type ITask interface {
	Do(i int)
}

type Task struct {
	fun func(i int)
}

func (t *Task) Do(i int) {
	t.fun(i)
}

type Consume struct {
	size     int
	len      int
	channels map[int]chan ITask
	isStop   bool
	wg       *sync.WaitGroup
}

func NewConsume(size, len int) *Consume {
	return &Consume{size, len, make(map[int]chan ITask), false, &sync.WaitGroup{}}
}

func (c *Consume) Start() {
	for i := 0; i < c.size; i++ {
		index := i

		c.channels[index] = make(chan ITask, c.len)

		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			for task := range c.channels[index] {
				task.Do(index)
			}
		}()
	}
}

func (c *Consume) StopTask() {
	c.isStop = true

	for _, task := range c.channels {
		close(task)
	}
}

func (c *Consume) AddTask(key int, task ITask) {
	if c.isStop {
		return
	}

	c.channels[c.findKey(key)] <- task
}

func (c *Consume) findKey(key int) int {
	return key % c.size
}

func (c *Consume) Wait() {
	c.wg.Wait()
}

func TaskFun(fun func(i int)) ITask {
	return &Task{fun}
}
