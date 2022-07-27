package handler

import (
	"fmt"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/logger"
)

// ExampleWebsocket 使用案例
type ExampleWebsocket struct {
}

func NewExampleWebsocket() *ExampleWebsocket {
	return &ExampleWebsocket{}
}

func (c *ExampleWebsocket) Connect(ctx *ichat.Context) error {
	conn, err := im.NewConnect(ctx.Context)
	if err != nil {
		logger.Error("websocket connect error: ", err.Error())
		return nil
	}

	// 创建客户端
	im.NewClient(ctx.RequestCtx(), conn, &im.ClientOptions{
		Channel: im.Session.Example,
		Uid:     0, // 自行提供用户ID
	}, im.NewClientCallback(
		// 连接成功回调
		im.WithOpenCallback(func(client im.IClient) {
			fmt.Printf("客户端[%d] 已连接\n", client.ClientId())
		}),
		// 接收消息回调
		im.WithMessageCallback(func(client im.IClient, message []byte) {
			// _ = message.Client.Write(&im.ClientOutContent{
			// 	IsAck:   true,
			// 	Content: []byte(message.Content),
			// }) // 推送消息
		}),
		// 关闭连接回调
		im.WithCloseCallback(func(client im.IClient, code int, text string) {
			fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.ClientId(), code, text)
		}),
	))

	return nil
}
