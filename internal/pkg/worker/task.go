package worker

import (
	"sync"
)

type Task struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func NewTask(buffer uint) *Task {
	if buffer == 0 {
		buffer = 1
	}

	return &Task{
		ch: make(chan struct{}, buffer),
		wg: &sync.WaitGroup{},
	}
}

func (w *Task) Do(fn func()) {

	if fn == nil {
		return
	}

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

func (w *Task) Wait() {
	w.wg.Wait()
	close(w.ch)
}
