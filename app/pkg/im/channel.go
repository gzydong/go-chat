package im

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/pkg/slice"
)

type HandleInterface interface {
	Open(client *Client)
	Message(message *ReceiveContent)
	Close(client *Client, code int, text string)
}

// ChannelManager 渠道管理（多渠道划分，实现不同业务之间隔离）
type ChannelManager struct {
	Name    string               // 渠道名称
	Count   int                  // 客户端连接数
	Clients map[int]*Client      // 客户端列表
	inChan  chan *ReceiveContent // 消息接收通道
	outChan chan *SenderContent  // 消息发送通道
	Lock    *sync.RWMutex        // 读写锁
	Handler HandleInterface      // 回调处理
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
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if _, ok := c.Clients[client.ClientId]; !ok {
		return false
	}

	delete(c.Clients, client.ClientId)

	c.Count--
	return true
}

// GetClient 获取客户端
func (c *ChannelManager) GetClient(cid int) (*Client, bool) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	client, ok := c.Clients[cid]

	return client, ok
}

// PushRecvChannel 推送消息到接收通道
func (c *ChannelManager) PushRecvChannel(message *ReceiveContent) {
	select {
	case c.inChan <- message:
		break
	case <-time.After(1000 * time.Millisecond):
		fmt.Printf("[%s] RecvChan 写入消息超时,管道长度：%d \n", c.Name, len(c.inChan))
		break
	}
}

// PushSendChannel 推送消息到消费通道
func (c *ChannelManager) PushSendChannel(msg *SenderContent) {
	select {
	case c.outChan <- msg:
		break
	case <-time.After(1000 * time.Millisecond):
		fmt.Printf("[%s] SendChan 写入消息超时,管道长度：%d \n", c.Name, len(c.inChan))
		break
	}
}

// SetCallbackHandler 设置 WebSocket 处理事件
func (c *ChannelManager) SetCallbackHandler(handle HandleInterface) *ChannelManager {
	c.Handler = handle

	return c
}

// Handle 渠道消费协程
func (c *ChannelManager) Handle(ctx context.Context) error {
	go c.recv(ctx)
	go c.send(ctx)

	return nil
}

// 接收客户端消息
func (c *ChannelManager) recv(ctx context.Context) {
	var (
		out     = 2 * time.Second
		timeout = time.NewTimer(out)
	)

	for {
		timeout.Reset(out)

		select {
		case <-ctx.Done():
			return

		// 处理接收消息
		case msg, ok := <-c.inChan:
			if ok {
				c.Handler.Message(msg)
			}

		case <-timeout.C:
		}
	}
}

// 推送客户端数据
func (c *ChannelManager) send(ctx context.Context) {
	var (
		out     = 2 * time.Second
		timeout = time.NewTimer(out)
	)

	for {
		timeout.Reset(out)

		select {
		case <-ctx.Done():
			return

		case body, ok := <-c.outChan:
			if ok {
				content, _ := jsonutil.JsonEncodeToByte(body.GetMessage())

				// 判断是否广播消息
				if body.IsBroadcast() {
					c.Lock.RLock()
					for cid, client := range c.Clients {
						if client.IsClosed || slice.InInt(cid, body.exclude) {
							continue
						}

						_ = client.Conn.WriteMessage(websocket.TextMessage, content)
					}
					c.Lock.RUnlock()
				} else {
					for _, cid := range body.receives {
						if client, ok := c.Clients[cid]; ok && !client.IsClosed {
							_ = client.Conn.WriteMessage(websocket.TextMessage, content)
						}
					}
				}
			}

		case <-timeout.C:
		}
	}
}
