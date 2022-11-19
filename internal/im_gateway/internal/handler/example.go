package handler

import (
	"go-chat/internal/im_gateway/internal/event"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

// ExampleChannel 使用案例
type ExampleChannel struct {
	storage *cache.ClientStorage
	event   *event.ExampleEvent
}

func NewExampleChannel(storage *cache.ClientStorage, event *event.ExampleEvent) *ExampleChannel {
	return &ExampleChannel{storage: storage, event: event}
}

func (c *ExampleChannel) WsConnect(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		logger.Errorf("websocket connect error: %s", err.Error())
		return nil
	}

	// 创建客户端
	im.NewClient(ctx.Ctx(), conn, &im.ClientOptions{
		Channel: im.Session.Example,
		Uid:     0, // 自行提供用户ID
	}, im.NewClientCallback(
		// 连接成功回调事件
		im.WithOpenCallback(c.event.OnOpen),
		// 接收消息回调
		im.WithMessageCallback(c.event.OnMessage),
		// 关闭连接回调
		im.WithCloseCallback(c.event.OnClose),
	))

	return nil
}
