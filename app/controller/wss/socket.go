package wss

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/pakg/im"
	"log"
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

	// 设置客户端主动断开连接触发事件
	client.SetCloseHandler(func(code int, text string) error {

		fmt.Println("客户端已关闭 ：", code, text)

		return nil
	})

	// 启动客户端心跳检测
	go client.Heartbeat(func(t *im.Client) bool {
		fmt.Printf("客户端心跳检测超时[%s]\n", t.Uuid)
		return true
	})

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

	// 设置客户端主动断开连接触发事件
	client.SetCloseHandler(func(code int, text string) error {

		fmt.Println("客户端已关闭 ：", code, text)

		return nil
	})

	// 启动客户端心跳检测
	go client.Heartbeat(func(t *im.Client) bool {
		fmt.Printf("客户端心跳检测超时[%s]\n", t.Uuid)
		return true
	})

	// 创建协程处理接收信息
	go client.AcceptClient()
}
