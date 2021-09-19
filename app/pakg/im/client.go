package im

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	heartbeatCheckInterval int = 20 // 心跳检测时间
	heartbeatIdleTime      int = 50 // 心跳超时时间
)

// Client WebSocket 客户端连接信息
type Client struct {
	Conn     *websocket.Conn // 客户端连接
	Uuid     string          // 客户端唯一标识
	UserId   int             // 用户ID
	LastTime int64           // 客户端最后心跳时间/心跳检测
	Channel  string          // 渠道分组
}

type CloseFunc func(c *Client) bool

// NewImClient ...
func NewImClient(conn *websocket.Conn, userId int) *Client {
	return &Client{
		Conn:     conn,
		Uuid:     uuid.NewV4().String(),
		UserId:   userId,
		LastTime: time.Now().Unix(),
	}
}

// Close 关闭客户端连接
func (w *Client) Close(code int, message string) {

	Handler := w.Conn.CloseHandler()
	_ = Handler(code, message)

	w.Conn.Close()
}

// Heartbeat 心跳检测
func (w *Client) Heartbeat(fn CloseFunc) {
	for {
		time.Sleep(time.Duration(heartbeatCheckInterval) * time.Second)

		if time.Now().Unix()-w.LastTime > int64(heartbeatIdleTime) {
			isOk := fn(w)
			if isOk {
				w.Close(500, "心跳检测超时，连接自动关闭")
				break
			}
		}
	}
}
