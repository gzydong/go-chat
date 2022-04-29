package process

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

var once sync.Once

type CoroutineInterface interface {
	Setup(ctx context.Context) error
}

type Coroutine struct {
	items []CoroutineInterface
}

func NewCoroutine(health *Health, subscribe *WsSubscribe) *Coroutine {
	coroutine := &Coroutine{}

	// 注册健康上报协程
	coroutine.register(health)

	// 注册消息订阅协程
	coroutine.register(subscribe)

	return coroutine
}

func (c *Coroutine) register(process CoroutineInterface) {
	c.items = append(c.items, process)
}

func (c *Coroutine) Start(eg *errgroup.Group, ctx context.Context) {
	once.Do(func() {
		for _, process := range c.items {
			func(obj CoroutineInterface) {
				eg.Go(func() error {
					return obj.Setup(ctx)
				})
			}(process)
		}
	})
}
