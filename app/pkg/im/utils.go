package im

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/app/helper"
	"net/http"
	"strconv"
	"time"
)

// NewWebsocket 获取 WebSocket 连接
func NewWebsocket(ctx *gin.Context) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	return upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
}

// NewClientID 创建客户端ID
func NewClientID() int {
	num := fmt.Sprintf("%03d", helper.MtRand(1, 999))

	val, _ := strconv.Atoi(fmt.Sprintf("%d%s", time.Now().UnixNano()/1000, num))

	return val
}
