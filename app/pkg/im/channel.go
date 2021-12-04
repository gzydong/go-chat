package im

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go-chat/app/pkg/jsonutil"
)

type HandleInterface interface {
	Open(client *Client)
	Message(message *ReceiveContent)
	Close(client *Client, code int, text string)
}

// Channel 渠道管理（多渠道划分，实现不同业务之间隔离）
type Channel struct {
	name    string               // 渠道名称
	count   int                  // 客户端连接数
	node    *Node                // 客户端列表【客户端ID取余拆分，降低 map 长度，减少 map 加锁时间提高并发处理量】
	inChan  chan *ReceiveContent // 消息接收通道
	outChan chan *SenderContent  // 消息发送通道
	handler HandleInterface      // 回调处理
}

func NewChannel(name string, node *Node, inChan chan *ReceiveContent, outChan chan *SenderContent) *Channel {
	return &Channel{name: name, node: node, inChan: inChan, outChan: outChan}
}

// Name 获取渠道名称
func (c *Channel) Name() string {
	return c.name
}

// Count 获取客户端连接数
func (c *Channel) Count() int {
	return c.count
}

// Client 获取客户端
func (c *Channel) Client(cid int64) (*Client, bool) {
	return c.node.get(cid)
}

// addClient 添加客户端
func (c *Channel) addClient(client *Client) {
	c.node.add(client)

	c.count++
}

// delClient 删除客户端
func (c *Channel) delClient(client *Client) {
	if !c.node.exist(client.cid) {
		return
	}

	c.node.del(client)
	c.count--
}

// PushRecvChannel 推送消息到接收通道
func (c *Channel) PushRecvChannel(message *ReceiveContent) {
	select {
	case c.inChan <- message:
		break
	case <-time.After(1000 * time.Millisecond):
		fmt.Printf("[%s] RecvChan 写入消息超时,管道长度：%d \n", c.name, len(c.inChan))
		break
	}
}

// PushSendChannel 推送消息到消费通道
func (c *Channel) PushSendChannel(msg *SenderContent) {
	select {
	case c.outChan <- msg:
		break
	case <-time.After(1000 * time.Millisecond):
		fmt.Printf("[%s] SendChan 写入消息超时,管道长度：%d \n", c.name, len(c.inChan))
		break
	}
}

// SetHandler 设置 WebSocket 处理事件
func (c *Channel) SetHandler(handle HandleInterface) *Channel {
	c.handler = handle

	return c
}

// Handle 渠道消费协程
func (c *Channel) Handle(ctx context.Context) error {
	go c.recv(ctx)
	go c.send(ctx)

	return nil
}

// 接收客户端消息
func (c *Channel) recv(ctx context.Context) {
	var (
		out     = 2 * time.Second
		timeout = time.NewTimer(out)
	)

	for {
		timeout.Reset(out)

		select {
		case <-ctx.Done():
			break

		// 处理接收消息
		case msg, ok := <-c.inChan:
			if ok {
				c.handler.Message(msg)
			}

		case <-timeout.C:
		}
	}
}

// 推送客户端数据
func (c *Channel) send(ctx context.Context) {
	var (
		out     = 2 * time.Second
		timeout = time.NewTimer(out)
	)

	for {
		timeout.Reset(out)

		select {
		case <-ctx.Done():
			break

		case body, ok := <-c.outChan:
			if ok {
				content, _ := jsonutil.JsonEncodeByte(body.GetMessage())

				// 判断是否广播消息
				if body.IsBroadcast() {
					c.node.each(func(c *Client) {
						_ = c.Write(websocket.TextMessage, content)
					})
				} else {
					for _, cid := range body.receives {
						if client, ok := c.Client(cid); ok {
							_ = client.Write(websocket.TextMessage, content)
						}
					}
				}
			}

		case <-timeout.C:
		}
	}
}
