package worker

import (
	"context"
	"sync"
	"sync/atomic"
)

type Worker struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	ch      chan func()
	isStop  bool
	counter int32
}

func NewWorker(worker int, buffer int) *Worker {

	if worker <= 0 {
		panic("worker must be greater than 0")
	}

	if buffer <= 0 {
		panic("buffer must be greater than 0")
	}

	c := &Worker{
		ch: make(chan func(), buffer),
		wg: &sync.WaitGroup{},
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.start(worker)

	return c
}

func (w *Worker) start(worker int) {

	w.wg.Add(worker)

	for i := 0; i < worker; i++ {
		go func() {
			defer w.wg.Done()

			for {
				select {
				case <-w.ctx.Done():
					return
				case task := <-w.ch:
					w.exec(task)
				}
			}
		}()
	}
}

func (w *Worker) exec(fn func()) {

	defer func() {
		if atomic.LoadInt32(&w.counter) == 0 && w.isStop {
			w.cancel()
		}
	}()

	fn()

	atomic.AddInt32(&w.counter, -1)
}

func (w *Worker) Do(fn func()) {

	if fn == nil {
		return
	}

	atomic.AddInt32(&w.counter, 1)

	w.ch <- fn
}

func (w *Worker) Wait() {
	w.isStop = true
	w.wg.Wait()
}
