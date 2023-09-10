package handler

import (
	"log"

	"go-chat/internal/commet/event"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/ichat/socket/adapter"
	"go-chat/internal/repository/cache"
)

type ChatChannel struct {
	Storage *cache.ClientStorage
	Event   *event.ChatEvent
}

// Conn 初始化连接
func (c *ChatChannel) Conn(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return c.NewClient(ctx.UserId(), conn)
}

func (c *ChatChannel) NewClient(uid int, conn socket.IConn) error {

	return socket.NewClient(conn, &socket.ClientOption{
		Uid:     uid,
		Channel: socket.Session.Chat,
		Storage: c.Storage,
		Buffer:  10,
	}, socket.NewEvent(
		// 连接成功回调
		socket.WithOpenEvent(c.Event.OnOpen),
		// 接收消息回调
		socket.WithMessageEvent(c.Event.OnMessage),
		// 关闭连接回调
		socket.WithCloseEvent(c.Event.OnClose),
	))
}
