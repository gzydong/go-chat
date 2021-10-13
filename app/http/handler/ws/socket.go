package ws

import (
	"go-chat/app/pkg/im"
	"log"

	"github.com/gin-gonic/gin"
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

	options := &im.ClientOption{
		Channel:       im.Manager.DefaultChannel,
		UserId:        c.GetInt("__user_id__"),
		ClientService: w.ClientService,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}

// AdminIo 连接客户端
func (w *WebSocket) AdminIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	options := &im.ClientOption{
		Channel:       im.Manager.AdminChannel,
		UserId:        c.GetInt("__user_id__"),
		ClientService: w.ClientService,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}
