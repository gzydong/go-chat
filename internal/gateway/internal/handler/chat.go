package handler

import (
	"context"
	"log"

	"go-chat/internal/gateway/internal/event"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/repository/cache"
)

type ChatChannel struct {
	storage *cache.ClientStorage
	event   *event.ChatEvent
}

func NewChatChannel(storage *cache.ClientStorage, event *event.ChatEvent) *ChatChannel {
	return &ChatChannel{storage: storage, event: event}
}

// Conn 初始化连接
func (c *ChatChannel) Conn(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return c.NewClient(ctx.Ctx(), ctx.UserId(), conn)
}

func (c *ChatChannel) NewClient(ctx context.Context, uid int, conn im.IConn) error {
	return im.NewClient(ctx, conn, &im.ClientOption{
		Uid:     uid,
		Channel: im.Session.Chat,
		Storage: c.storage,
		Buffer:  10,
	}, im.NewClientCallback(
		// 连接成功回调事件
		im.WithOpenCallback(c.event.OnOpen),
		// 接收消息回调
		im.WithMessageCallback(c.event.OnMessage),
		// 关闭连接回调
		im.WithCloseCallback(c.event.OnClose),
	))
}
