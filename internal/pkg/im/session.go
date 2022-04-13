package im

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
	Default *Channel // 默认分组
	Example *Channel // 案例分组

	// 可自行注册其它渠道...
}

func Initialize(ctx context.Context, eg *errgroup.Group) {
	once.Do(func() {
		initialize(ctx, eg)
	})
}

func initialize(ctx context.Context, eg *errgroup.Group) {
	Session = &session{
		Default: NewChannel("default", NewNode(10), make(chan *SenderContent, 5<<20)),
		Example: NewChannel("example", NewNode(1), make(chan *SenderContent, 100)),
	}

	// 延时启动守护协程
	time.AfterFunc(5*time.Second, func() {
		eg.Go(func() error {
			return health.Start(ctx)
		})

		eg.Go(func() error {
			return ack.Start(ctx)
		})

		eg.Go(func() error {
			return Session.Default.Start(ctx)
		})

		eg.Go(func() error {
			return Session.Example.Start(ctx)
		})
	})
}
