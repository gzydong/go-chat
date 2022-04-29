package worker

import (
	"sync"
)

type Concurrent struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func NewConcurrent(num int) *Concurrent {
	if num == 0 {
		num = 1
	}

	return &Concurrent{
		ch: make(chan struct{}, num),
		wg: &sync.WaitGroup{},
	}
}

func (w *Concurrent) Add(fn func()) {
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

type TaskInterface interface {
	Do()
}

func (w *Concurrent) AddTask(task TaskInterface) {
	w.Add(task.Do)
}

func (w *Concurrent) Wait() {
	w.wg.Wait()
	close(w.ch)
}
