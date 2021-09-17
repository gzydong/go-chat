package socket

import (
	"sync"
)

// ChannelManager WebSocket 渠道管理
type ChannelManager struct {
	ChannelName string               // 渠道名称
	ConnectNum  int                  // 渠道客户端连接数
	Clients     map[string]*WsClient // 客户端列表
	Lock        *sync.Mutex
}

// PushClient
func (c *ChannelManager) PushClient(client *WsClient) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Clients[client.Uuid] = client

	c.ConnectNum++
}

// RemoveClient 移出客户端
func (c *ChannelManager) RemoveClient(client *WsClient) bool {
	_, ok := c.Clients[client.Uuid]
	if !ok {
		return false
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Clients, client.Uuid)

	c.ConnectNum--

	return true
}
