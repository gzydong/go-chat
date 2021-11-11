package im

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// NewWebsocket 获取 WebSocket 连接
func NewWebsocket(ctx *gin.Context) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024 * 2, // 指定读缓存区大小
		WriteBufferSize: 1024 * 2, // 指定写缓存区大小
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	return upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
}
