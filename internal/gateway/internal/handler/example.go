package handler

import (
	"log"

	"go-chat/internal/gateway/internal/event"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/repository/cache"
)

// ExampleChannel 案例
type ExampleChannel struct {
	storage *cache.ClientStorage
	event   *event.ExampleEvent
}

func NewExampleChannel(storage *cache.ClientStorage, event *event.ExampleEvent) *ExampleChannel {
	return &ExampleChannel{storage: storage, event: event}
}

func (c *ExampleChannel) Conn(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return im.NewClient(ctx.Ctx(), conn, &im.ClientOption{
		Channel: im.Session.Example,
		Uid:     0,
	}, im.NewClientCallback(
		// 连接成功回调事件
		im.WithOpenCallback(c.event.OnOpen),
		// 接收消息回调
		im.WithMessageCallback(c.event.OnMessage),
		// 关闭连接回调
		im.WithCloseCallback(c.event.OnClose),
	))
}
