package socket

import (
	"github.com/gorilla/websocket"
	"sync"
)

// WsClient WebSocket 客户端连接信息
type WsClient struct {
	Conn     *websocket.Conn // 客户端连接
	Uuid     string          // 客户端唯一标识
	UserId   int             // 用户ID
	LastTime int64           // 客户端最后心跳时间
}

// 渠道客户端
type ChannelClient struct {
	DefaultChannel *ChannelManager
}

var Manager *ChannelClient

func init() {
	Manager = &ChannelClient{
		DefaultChannel: &ChannelManager{
			ChannelName: "default",
			ConnectNum:  0,
			Clients:     make(map[string]*WsClient),
			Lock:        &sync.Mutex{},
		},
	}
}
