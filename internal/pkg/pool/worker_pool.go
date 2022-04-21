package pool

import (
	"sync"
)

type WorkerPool struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

type taskInterface interface {
	Do()
}

func NewWorkerPool(num int) *WorkerPool {
	if num == 0 {
		num = 1
	}

	return &WorkerPool{
		ch: make(chan struct{}, num),
		wg: &sync.WaitGroup{},
	}
}

func (w *WorkerPool) Add(fn func()) {
	w.ch <- struct{}{}
	w.wg.Add(1)
	go func() {
		defer func() {
			w.wg.Done()
			<-w.ch
		}()

		fn()
	}()
}

func (w *WorkerPool) AddTask(task taskInterface) {
	w.Add(func() {
		task.Do()
	})
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}
