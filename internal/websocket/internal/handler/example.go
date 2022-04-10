package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jwt"
)

// ExampleWebsocket 使用案例
type ExampleWebsocket struct {
}

func (c *ExampleWebsocket) Connect(ctx *gin.Context) {
	conn, err := im.NewWebsocket(ctx)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	// 创建客户端
	im.NewClient(conn, &im.ClientOptions{
		Channel: im.Sessions.Default,
		Uid:     jwt.GetUid(ctx),
		CallBack: im.NewClientCallBack(im.WithClientCallBackOpen(func(client im.ClientInterface) {
			fmt.Printf("客户端[%d] 已连接", client.ClientId())
		}), im.WithClientCallBackMessage(func(message *im.ReceiveContent) {
			// 推送消息
			_ = message.Client.Write([]byte(message.Content))
		}), im.WithClientCallBackClose(func(client im.ClientInterface, code int, text string) {
			fmt.Printf("客户端[%d] 已关闭连接", client.ClientId())
		})),
	})
}
