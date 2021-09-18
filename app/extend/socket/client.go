package socket

import "github.com/gorilla/websocket"

// Client WebSocket 客户端连接信息
type Client struct {
	Conn     *websocket.Conn // 客户端连接
	Uuid     string          // 客户端唯一标识
	UserId   int             // 用户ID
	LastTime int64           // 客户端最后心跳时间
	Channel  string          // 渠道分组名
}

// Close 关闭客户端连接
func (w *Client) Close(code int, message string) {
	w.Conn.Close()
}
