package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
	"log"
)

type DefaultWebSocket struct {
	client *service.ClientService
}

func NewDefaultWebSocket(client *service.ClientService) *DefaultWebSocket {
	handler := &DefaultWebSocket{client: client}

	channel := im.Manager.DefaultChannel

	channel.SetCallbackHandler(handler)

	return &DefaultWebSocket{client: client}
}

// Connect 初始化连接
func (ws *DefaultWebSocket) Connect(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	options := &im.ClientOption{
		Channel:       im.Manager.DefaultChannel,
		UserId:        auth.GetAuthUserID(c),
		ClientService: ws.client,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}

// Open 连接成功回调事件
func (ws *DefaultWebSocket) Open(client *im.Client) {
	// fmt.Printf("[%s] 客户端已连接[%d] \n", client.Channel.Name, client.ClientId)
}

// Message 消息接收回调事件
func (ws *DefaultWebSocket) Message(message *im.ClientContent) {
	fmt.Printf("[%s]消息通知 Client:%d，Content: %s \n", message.Client.Channel.Name, message.Client.ClientId, message.Content)

	body := im.NewSenderContent().SetBroadcast(true).SetMessage(&im.Message{
		Event: "test",
		Content: &map[string]interface{}{
			"name":     "anskjfna",
			"nickname": "那可就散你氨基酸卡那",
		},
	})

	im.Manager.DefaultChannel.SendMessage(body)
}

// Close 客户端关闭回调事件
func (ws *DefaultWebSocket) Close(client *im.Client, code int, text string) {
	// fmt.Printf("[%s] 客户端[%d] 已关闭 code:%d text:%s \n", client.Channel.Name, client.ClientId, code, text)
}
