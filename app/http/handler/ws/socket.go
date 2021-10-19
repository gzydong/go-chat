package ws

import (
	"go-chat/app/pkg/auth"
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
)

type WebSocket struct {
	clientService *service.ClientService
}

func NewWebSocketHandler(client *service.ClientService) *WebSocket {
	return &WebSocket{clientService: client}
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
		UserId:        auth.GetAuthUserID(c),
		ClientService: w.clientService,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}
