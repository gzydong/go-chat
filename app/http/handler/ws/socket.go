package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
	"log"
)

type WebSocketTalk struct {
	client *service.ClientService
}

func NewWebSocketTalkHandler(client *service.ClientService) *WebSocketTalk {
	handler := &WebSocketTalk{client: client}

	channel := im.Manager.DefaultChannel

	channel.SetCallbackHandler(handler)

	return handler
}

func (w *WebSocketTalk) Connect(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	options := &im.ClientOption{
		Channel:       im.Manager.DefaultChannel,
		UserId:        auth.GetAuthUserID(c),
		ClientService: w.client,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}

// Open 连接成功回调事件
func (w *WebSocketTalk) Open(client *im.Client) {
	fmt.Printf("[%s] 客户端已连接[%d] \n", client.Channel.Name, client.ClientId)
}

// Message 消息接收回调事件
func (w *WebSocketTalk) Message(message *im.RecvMessage) {
	fmt.Printf("[%s]消息通知 Client:%d，Content: %s \n", message.Client.Channel.Name, message.Client.ClientId, message.Content)

	message.Client.Channel.SendMessage(&im.SendMessage{
		IsAll:   true,
		Clients: nil,
		Event:   "talk",
		Content: message.Content,
	})
}

// Close 客户端关闭回调事件
func (w *WebSocketTalk) Close(client *im.Client, code int, text string) {
	fmt.Printf("[%s] 客户端[%d] 已关闭 code:%d text:%s \n", client.Channel.Name, client.ClientId, code, text)
}
