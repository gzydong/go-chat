package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/im/adapter"
)

// ExampleWebsocket 使用案例
type ExampleWebsocket struct {
}

func NewExampleWebsocket() *ExampleWebsocket {
	return &ExampleWebsocket{}
}

func (c *ExampleWebsocket) Connect(ctx *ichat.Context) error {
	conn, err := adapter.NewWsAdapter(ctx.Context)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", err.Error())
		return nil
	}

	// 创建客户端
	im.NewClient(ctx.RequestCtx(), conn, &im.ClientOptions{
		Channel: im.Session.Example,
		Uid:     0, // 自行提供用户ID
	}, im.NewClientCallback(
		// 连接成功回调
		im.WithOpenCallback(func(client im.IClient) {
			fmt.Printf("客户端[%d] 已连接\n", client.Cid())
		}),
		// 接收消息回调
		im.WithMessageCallback(func(client im.IClient, message []byte) {
			fmt.Println("接收消息===>>>", message)
		}),
		// 关闭连接回调
		im.WithCloseCallback(func(client im.IClient, code int, text string) {
			fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.Cid(), code, text)
		}),
	))

	return nil
}
