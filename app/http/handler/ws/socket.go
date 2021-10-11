package ws

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
)

type WebSocket struct {
	ClientService *service.ClientService
}

// SocketIo 连接客户端
func (w *WebSocket) SocketIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	// 创建客户端
	im.NewClient(conn, w.ClientService, c.GetInt("__user_id__"), im.Manager.DefaultChannel).Start()
}

// AdminIo 连接客户端
func (w *WebSocket) AdminIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	// 创建客户端
	im.NewClient(conn, w.ClientService, c.GetInt("__user_id__"), im.Manager.AdminChannel).Start()
}
