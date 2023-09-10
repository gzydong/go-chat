package handler

import (
	"log"

	"go-chat/internal/commet/event"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/ichat/socket/adapter"
	"go-chat/internal/repository/cache"
)

// ExampleChannel 案例
type ExampleChannel struct {
	Storage *cache.ClientStorage
	Event   *event.ExampleEvent
}

func (c *ExampleChannel) Conn(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return socket.NewClient(conn, &socket.ClientOption{
		Channel: socket.Session.Example,
		Uid:     0,
	}, socket.NewEvent(
		// 连接成功回调
		socket.WithOpenEvent(c.Event.OnOpen),
		// 接收消息回调
		socket.WithMessageEvent(c.Event.OnMessage),
		// 关闭连接回调
		socket.WithCloseEvent(c.Event.OnClose),
	))
}
