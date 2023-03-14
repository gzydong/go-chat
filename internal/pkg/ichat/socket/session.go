package socket

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Session 客户端管理实例
var Session *session

var once sync.Once

// session 渠道客户端
type session struct {
	Chat    *Channel // 默认分组
	Example *Channel // 案例分组

	// 可自行注册其它渠道...
}

func Initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	once.Do(func() {
		initialize(ctx, eg, fn)
	})
}

func initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	Session = &session{
		Chat:    NewChannel("chat", make(chan *SenderContent, 5<<20)),
		Example: NewChannel("example", make(chan *SenderContent, 100)),
	}

	// 延时启动守护协程
	time.AfterFunc(5*time.Second, func() {
		eg.Go(func() error {
			defer fn("health exit")
			return health.Start(ctx)
		})

		eg.Go(func() error {
			defer fn("ack exit")
			return ack.Start(ctx)
		})

		eg.Go(func() error {
			defer fn("chat exit")
			return Session.Chat.Start(ctx)
		})

		eg.Go(func() error {
			defer fn("example exit")
			return Session.Example.Start(ctx)
		})
	})
}
