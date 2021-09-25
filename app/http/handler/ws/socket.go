package ws

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/im"
)

type WsController struct {
}

// SocketIo 连接客户端
func (w *WsController) SocketIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error")
		return
	}

	// 创建客户端
	client := im.NewImClient(conn, c.GetInt("user_id"), im.Manager.DefaultChannel)

	// 启动客户端心跳检测
	go client.Heartbeat()

	// 创建协程处理接收信息
	go client.AcceptClient()
}

// AdminIo 连接客户端
func (w *WsController) AdminIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error")
		return
	}

	// 创建客户端
	client := im.NewImClient(conn, c.GetInt("user_id"), im.Manager.AdminChannel)

	// 启动客户端心跳检测
	go client.Heartbeat()

	// 创建协程处理接收信息
	go client.AcceptClient()
}
