package im

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

type IChannel interface {
	Name() string
	Count() int64
}

// Channel 渠道管理（多渠道划分，实现不同业务之间隔离）
type Channel struct {
	name    string              // 渠道名称
	count   int64               // 客户端连接数
	node    *Node               // 客户端列表【客户端ID取余拆分，降低 map 长度】
	outChan chan *SenderContent // 消息发送通道
}

func NewChannel(name string, node *Node, outChan chan *SenderContent) *Channel {
	return &Channel{name: name, node: node, outChan: outChan}
}

// Name 获取渠道名称
func (c *Channel) Name() string {
	return c.name
}

// Count 获取客户端连接数
func (c *Channel) Count() int64 {
	return c.count
}

// Client 获取客户端
func (c *Channel) Client(cid int64) (*Client, bool) {
	return c.node.get(cid)
}

// Write 推送消息到消费通道
func (c *Channel) Write(msg *SenderContent) {
	select {
	case c.outChan <- msg:
		break
	case <-time.After(3 * time.Second):
		fmt.Printf("[%s] Channel OutChan 写入消息超时,管道长度：%d \n", c.name, len(c.outChan))
		break
	}
}

// addClient 添加客户端
func (c *Channel) addClient(client *Client) {
	c.node.add(client)

	atomic.AddInt64(&c.count, 1)
}

// delClient 删除客户端
func (c *Channel) delClient(client *Client) {
	if !c.node.exist(client.cid) {
		return
	}

	c.node.del(client)

	atomic.AddInt64(&c.count, -1)
}

// 推送客户端数据
func (c *Channel) loopPush(ctx context.Context) {

	out := 2 * time.Second
	timer := time.NewTimer(out)

	defer timer.Stop()

	for {
		timer.Reset(out)

		select {
		case <-timer.C:
		case <-ctx.Done():
			return
		case body, ok := <-c.outChan:
			if ok {
				content, _ := json.Marshal(body.GetMessage())

				// 判断是否广播消息
				if body.IsBroadcast() {
					c.node.each(func(client *Client) {
						_ = client.Write(&ClientOutContent{Content: content})
					})
				} else {
					for _, cid := range body.receives {
						if client, ok := c.Client(cid); ok {
							_ = client.Write(&ClientOutContent{Content: content})
						}
					}
				}
			}
		}
	}
}

// Start 渠道消费协程
func (c *Channel) Start(ctx context.Context) error {
	go c.loopPush(ctx)

	return nil
}
