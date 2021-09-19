package wss

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/app/pakg/im"
	"log"
	"net/http"
	"time"
)

type WsController struct {
}

// SocketIo 连接客户端
func (w *WsController) SocketIo(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket connect error")
		return
	}

	client := im.NewImClient(conn, c.GetInt("user_id"))

	fmt.Printf("UserID: %#v, UUID: %s\n", client.UserId, client.Uuid)

	// 设置客户端主动断开连接触发事件
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("客户端已关闭 ：", code, text)
		im.Manager.DefaultChannel.RemoveClient(client)

		_ = conn.Close()
		return nil
	})

	im.Manager.DefaultChannel.RegisterClient(client)

	// 启动客户端心跳检测
	go client.Heartbeat(func(t *im.Client) bool {
		fmt.Printf("客户端心跳检测超时[%s]\n", t.Uuid)

		return true
	})

	go recv(client)
}

// recv 消息接收处理
func recv(client *im.Client) {
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

		im.Manager.DefaultChannel.SendMessage(&im.Message{
			Receiver: []string{client.Uuid},
			IsAll:    false,
			Event:    "talk_type",
			Content:  string(message),
		})
	}
}
