package im

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type WebsocketInterface interface {
	Open(client *Client)
	Message(message *RecvMessage)
	Close(client *Client, code int, text string)
}

// ChannelManager WebSocket 渠道管理（多渠道划分，实现不同业务之间隔离）
type ChannelManager struct {
	Name     string            // 渠道名称
	Count    int               // 客户端连接数
	Clients  map[int]*Client   // 客户端列表
	RecvChan chan *RecvMessage // 消息接收通道
	SendChan chan *SendMessage // 消息发送通道
	Lock     *sync.Mutex       // 互斥锁
	Handle   WebsocketInterface
}

// RegisterClient 注册客户端
func (c *ChannelManager) RegisterClient(client *Client) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Clients[client.ClientId] = client

	c.Count++
}

// RemoveClient 删除客户端
func (c *ChannelManager) RemoveClient(client *Client) bool {
	_, ok := c.Clients[client.ClientId]
	if !ok {
		return false
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Clients, client.ClientId)

	c.Count--
	return true
}

// GetClient 获取客户端
func (c *ChannelManager) GetClient(clientId int) (*Client, bool) {
	client, ok := c.Clients[clientId]

	return client, ok
}

// RecvMessage 推送消息到消息接收通道
func (c *ChannelManager) RecvMessage(message *RecvMessage) {
	select {
	case c.RecvChan <- message:
		break
	case <-time.After(800 * time.Millisecond):
		fmt.Printf("[%s] RecvChan 写入消息超时,管道长度：%d \n", c.Name, len(c.RecvChan))
		break
	}
}

// SendMessage 推送消息到消费通道
func (c *ChannelManager) SendMessage(message *SendMessage) {
	select {
	case c.SendChan <- message:
		break
	case <-time.After(800 * time.Millisecond):
		fmt.Printf("[%s] SendChan 写入消息超时,管道长度：%d \n", c.Name, len(c.RecvChan))
		break
	}
}

// SetCallbackHandler 设置 WebSocket 处理事件
func (c *ChannelManager) SetCallbackHandler(handle WebsocketInterface) *ChannelManager {
	c.Handle = handle

	return c
}

// Process 渠道消费协程
func (c *ChannelManager) Process(ctx context.Context) {
	go c.RecvProcess(ctx)
	go c.SendProcess(ctx)
}

func (c *ChannelManager) RecvProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		// 处理接收消息
		case value, ok := <-c.RecvChan:
			if ok {
				c.Handle.Message(value)
			}

			break
		case <-time.After(3 * time.Second):
			break
		}
	}
}

func (c *ChannelManager) SendProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-c.SendChan:
			if !ok {
				break
			}

			content, _ := jsoniter.Marshal(value)

			// 判断是否推送所有客户端
			if value.IsAll {
				c.Lock.Lock() // 保证线程安全
				for _, client := range c.Clients {
					if client.IsClose {
						continue
					}

					_ = client.Conn.WriteMessage(websocket.TextMessage, content)
				}
				c.Lock.Unlock()
			} else {
				for _, clientId := range value.Clients {
					client, ok := c.Clients[clientId]
					if ok && client.IsClose == false {
						_ = client.Conn.WriteMessage(websocket.TextMessage, content)
					}
				}
			}

			break
		case <-time.After(3 * time.Second):

			break
		}
	}
}
