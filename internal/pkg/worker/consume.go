package worker

import (
	"sync"
)

type ITask interface {
	Do(workerId int)
}

type Task struct {
	do func(workerId int)
}

func (t *Task) Do(workerId int) {
	t.do(workerId)
}

type Consume struct {
	worker   int
	len      int
	channels map[int]chan ITask
	isStop   bool
	wg       *sync.WaitGroup
}

func NewConsume(worker, len int) *Consume {

	consume := &Consume{
		worker:   worker,
		len:      len,
		channels: make(map[int]chan ITask),
		wg:       &sync.WaitGroup{},
	}

	consume.Start()

	return consume
}

func (c *Consume) Start() {
	for i := 0; i < c.worker; i++ {
		index := i

		task := make(chan ITask, c.len)

		c.channels[index] = task

		c.wg.Add(1)

		go func() {
			defer c.wg.Done()

			for iTask := range task {
				iTask.Do(index)
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

	index := c.index(key)

	c.channels[index] <- task
}

func (c *Consume) index(key int) int {
	return key % c.worker
}

func (c *Consume) Wait() {
	c.wg.Wait()
}

func TaskFunc(do func(i int)) ITask {
	return &Task{do}
}
