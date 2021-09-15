package wss

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WsController struct {
}

// WsClient 连接客户端
func (w *WsController) WsClient(c *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", c.Param("channel"))
		return
	}

	go recv(conn)
}

func recv(conn *websocket.Conn) {
	defer conn.Close()

	// 最后一次发送消息的时间
	lastTime := time.Now().Unix()

	// 创建一个周期性的定时器,用做心跳检测
	ticker := time.NewTicker(20 * time.Second)
	go func(conn *websocket.Conn) {
		for {
			<-ticker.C

			if time.Now().Unix()-lastTime > 50 {
				ticker.Stop()
				conn.Close()
				fmt.Println("心跳检测超时，连接自动关闭")
			}
		}
	}(conn)

	for {
		//读取ws中的数据
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		lastTime = time.Now().Unix()

		if string(message) == "ping" {
			message = []byte("pong")
		}

		//写入ws数据
		err = conn.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
