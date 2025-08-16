package process

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type WorkerGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	sync.WaitGroup
}

func NewWorkerGroup(ctx context.Context) *WorkerGroup {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		// 创建一个通道来接收信号
		sigChan := make(chan os.Signal, 1)

		// 监听指定的信号（例如 syscall.SIGINT 和 syscall.SIGTERM）
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-sigChan:
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}()

	return &WorkerGroup{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *WorkerGroup) Go(fn func(ctx context.Context)) {
	w.Add(1)
	go func() {
		defer w.Done()
		fn(w.ctx)
	}()
}
