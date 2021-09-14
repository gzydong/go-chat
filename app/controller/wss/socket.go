package wss

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

	for {
		//读取ws中的数据
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

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
