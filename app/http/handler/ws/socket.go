package ws

import (
	"go-chat/app/pkg/auth"
	"go-chat/app/websocket"
	"log"

	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
)

type WebSocket struct {
	client  *service.ClientService
	channel *im.ChannelManager
}

func NewWebSocketHandler(client *service.ClientService) *WebSocket {

	channel := im.Manager.DefaultChannel

	channel.SetCallbackHandler(websocket.NewDefaultChannelHandle())

	return &WebSocket{client: client, channel: channel}
}

// SocketIo 连接客户端
func (w *WebSocket) SocketIo(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	options := &im.ClientOption{
		Channel:       w.channel,
		UserId:        auth.GetAuthUserID(c),
		ClientService: w.client,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}
