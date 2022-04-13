package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/im"
)

// ExampleWebsocket 使用案例
type ExampleWebsocket struct {
}

func NewExampleWebsocket() *ExampleWebsocket {
	return &ExampleWebsocket{}
}

func (c *ExampleWebsocket) Connect(ctx *gin.Context) {
	conn, err := im.NewWebsocket(ctx)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	// 创建客户端
	im.NewClient(conn, &im.ClientOptions{
		Channel: im.Session.Example,
		Uid:     0, // 自行提供用户ID
	}, im.NewClientCallback(
		// 连接成功回调
		im.WithOpenCallback(func(client im.ClientInterface) {
			fmt.Printf("客户端[%d] 已连接\n", client.ClientId())
		}),
		// 接收消息回调
		im.WithMessageCallback(func(message *im.ReceiveContent) {
			_ = message.Client.Write(&im.ClientOutContent{
				IsAck:   true,
				Content: []byte(message.Content),
			}) // 推送消息
		}),
		// 关闭连接回调
		im.WithCloseCallback(func(client im.ClientInterface, code int, text string) {
			fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.ClientId(), code, text)
		}),
	))
}
