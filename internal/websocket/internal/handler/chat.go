package handler

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
	"go-chat/internal/websocket/internal/event"
)

type ChatChannel struct {
	storage *cache.ClientStorage
	event   *event.ChatEvent
}

func NewChatChannel(storage *cache.ClientStorage, event *event.ChatEvent) *ChatChannel {
	return &ChatChannel{storage: storage, event: event}
}

// WsConn 初始化连接
func (c *ChatChannel) WsConn(ctx *ichat.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		logger.Errorf("websocket connect error: %s", err.Error())
		return nil
	}

	// 创建客户端
	c.client(ctx.Ctx(), ctx.UserId(), conn)

	return nil
}

// TcpConn 初始化连接
func (c *ChatChannel) TcpConn(ctx context.Context, conn *adapter.TcpAdapter) {
	c.client(ctx, 2054, conn)
}

func (c *ChatChannel) client(ctx context.Context, uid int, conn im.IConn) {
	im.NewClient(ctx, conn, &im.ClientOptions{
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
