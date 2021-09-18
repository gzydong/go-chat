package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// ChannelManager WebSocket 渠道管理（多渠道划分，实现不同业务之间隔离）
type ChannelManager struct {
	Name        string             // 渠道名称
	Count       int                // 渠道客户端连接数统计
	Clients     map[string]*Client // 客户端列表
	ChanMessage chan *Message      // 消息发送通道
	Lock        *sync.Mutex
}

// RegisterClient 注册客户端
func (c *ChannelManager) RegisterClient(client *Client) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// 设置渠道名称
	client.Channel = c.Name

	c.Clients[client.Uuid] = client

	c.Count++
}

// RemoveClient 移出客户端
func (c *ChannelManager) RemoveClient(client *Client) bool {
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

// GetClient 获取指定客户端
func (c *ChannelManager) GetClient(uuid string) (*Client, bool) {
	client, ok := c.Clients[uuid]

	return client, ok
}

// ConsumerProcess 渠道消费协程
func (c *ChannelManager) ConsumerProcess() {
	fmt.Println("消费协程已启动: ", c.Name)

	for {
		select {
		case value, ok := <-c.ChanMessage:
			if !ok {
				fmt.Printf("消费通道[%s]，读取数据失败...", c.Name)
				return
			}

			content, _ := json.Marshal(value)

			// 判断是否推送所有客户端
			if value.IsAll {
				for _, client := range c.Clients {
					_ = client.Conn.WriteMessage(websocket.TextMessage, content)
				}
			} else {
				for _, uuid := range value.Receiver {
					client, ok := c.Clients[uuid]
					if ok {
						_ = client.Conn.WriteMessage(websocket.TextMessage, content)
					}
				}
			}
			break
		}
	}
}
