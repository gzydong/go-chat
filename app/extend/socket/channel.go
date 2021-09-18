package socket

import (
	"sync"
)

// ChannelManager WebSocket 渠道管理（多渠道划分，实现不同业务之间隔离）
type ChannelManager struct {
	ChannelName string               // 渠道名称
	Count       int                  // 渠道客户端连接数统计
	Clients     map[string]*WsClient // 客户端列表
	Lock        *sync.Mutex
}

// RegisterClient 注册客户端
func (c *ChannelManager) RegisterClient(client *WsClient) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// 设置渠道名称
	client.Channel = c.ChannelName

	c.Clients[client.Uuid] = client

	c.Count++
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

	c.Count--

	return true
}
