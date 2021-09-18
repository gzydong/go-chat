package wss

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"go-chat/app/extend/socket"
	"log"
	"net/http"
	"time"
)

type WsController struct {
}

// WsClient 连接客户端
func (w *WsController) WsClient(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", c.Param("channel"))
		return
	}

	client := &socket.WsClient{
		Conn:     conn,
		Uuid:     uuid.NewV4().String(),
		UserId:   c.GetInt("user_id"),
		LastTime: time.Now().Unix(),
	}

	socket.Manager.DefaultChannel.RegisterClient(client)

	// 设置客户端主动断开连接触发事件
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("客户端已关闭 ：", code, text)
		socket.Manager.DefaultChannel.RemoveClient(client)

		_ = conn.Close()
		return nil
	})

	go heartbeat(client)
	go recv(client)
}

// heartbeat 心跳检测
func heartbeat(client *socket.WsClient) {
	// 创建一个周期性的定时器,用做心跳检测
	ticker := time.NewTicker(20 * time.Second)

	for {
		<-ticker.C

		if time.Now().Unix()-client.LastTime > 50 {
			ticker.Stop()

			Handler := client.Conn.CloseHandler()
			_ = Handler(500, "心跳检测超时，连接自动关闭")
		}
	}
}

// recv 消息接收处理
func recv(client *socket.WsClient) {
	defer client.Conn.Close()

	for {
		//读取ws中的数据
		mt, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 更新最后一次接受消息时间，用做心跳检测判断
		client.LastTime = time.Now().Unix()

		if string(message) == "ping" {
			message = []byte("pong")

			//写入ws数据
			err = client.Conn.WriteMessage(mt, message)
			if err != nil {
				break
			}

			continue
		}
	}
}
