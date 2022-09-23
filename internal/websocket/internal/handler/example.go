package handler

import (
	"github.com/sirupsen/logrus"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/event"
)

// ExampleChannel 使用案例
type ExampleChannel struct {
	cache *service.ClientService
	event *event.ExampleEvent
}

func NewExampleChannel(cache *service.ClientService, event *event.ExampleEvent) *ExampleChannel {
	return &ExampleChannel{cache: cache, event: event}
}

func (c *ExampleChannel) WsConnect(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", err.Error())
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
