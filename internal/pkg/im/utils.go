package im

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NewConnect 获取 WebSocket 连接
func NewConnect(ctx *gin.Context) (*websocket.Conn, error) {
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

// 获取 sync.Map 切片
func maps(num int) []*sync.Map {
	if num <= 0 {
		num = 1
	}

	items := make([]*sync.Map, 0, num)

	for i := 0; i < num; i++ {
		items = append(items, &sync.Map{})
	}

	return items
}
