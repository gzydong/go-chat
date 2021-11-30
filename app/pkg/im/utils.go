package im

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
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

// 拆分数量
func maps(num int) []*sync.Map {
	var items []*sync.Map

	if num <= 0 {
		num = 1
	}

	for i := 0; i < num; i++ {
		items = append(items, &sync.Map{})
	}

	return items
}

// 获取客户端ID在第几个 map 中
func getMapIndex(cid int64, num int) int {
	return int(cid) % num
}
