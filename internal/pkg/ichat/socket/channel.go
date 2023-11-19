package socket

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sourcegraph/conc/pool"
)

type IChannel interface {
	Name() string
	Count() int64
	Client(cid int64) (*Client, bool)
	Write(data *SenderContent)
	addClient(client *Client)
	delClient(client *Client)
}

// Channel 渠道管理（多渠道划分，实现不同业务之间隔离）
type Channel struct {
	name    string                              // 渠道名称
	count   int64                               // 客户端连接数
	node    cmap.ConcurrentMap[string, *Client] // 客户端列表
	outChan chan *SenderContent                 // 消息发送通道
}

func NewChannel(name string, outChan chan *SenderContent) *Channel {
	return &Channel{name: name, node: cmap.New[*Client](), outChan: outChan}
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
	return c.node.Get(strconv.FormatInt(cid, 10))
}

// Write 推送消息到消费通道
func (c *Channel) Write(data *SenderContent) {

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case c.outChan <- data:
	case <-timer.C:
		log.Printf("[ERROR] [%s] Channel OutChan 写入消息超时,管道长度：%d \n", c.name, len(c.outChan))
	}
}

// Start 渠道消费开启
func (c *Channel) Start(ctx context.Context) error {

	var (
		worker = pool.New().WithMaxGoroutines(10)
		timer  = time.NewTicker(15 * time.Second)
	)

	defer log.Println(fmt.Errorf("channel exit :%s", c.Name()))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("channel exit :%s", c.Name())
		case <-timer.C:
			// fmt.Printf("channel empty message name:%s unix:%d len:%d\n", c.name, time.Now().Unix(), len(c.outChan))
		case val, ok := <-c.outChan:
			if !ok {
				return fmt.Errorf("outchan close :%s", c.Name())
			}

			c.consume(worker, val, func(data *SenderContent, value *Client) {
				_ = value.Write(&ClientResponse{
					IsAck:   data.IsAck,
					Event:   data.message.Event,
					Content: data.message.Content,
					Retry:   3,
				})
			})
		}
	}
}

func (c *Channel) consume(worker *pool.Pool, data *SenderContent, fn func(data *SenderContent, value *Client)) {
	worker.Go(func() {

		if data.IsBroadcast() {
			c.node.IterCb(func(_ string, client *Client) {
				fn(data, client)
			})
			return
		}

		for _, cid := range data.receives {
			if client, ok := c.Client(cid); ok {
				fn(data, client)
			}
		}
	})
}

// addClient 添加客户端
func (c *Channel) addClient(client *Client) {
	c.node.Set(strconv.FormatInt(client.cid, 10), client)

	atomic.AddInt64(&c.count, 1)
}

// delClient 删除客户端
func (c *Channel) delClient(client *Client) {

	cid := strconv.FormatInt(client.cid, 10)

	if !c.node.Has(cid) {
		return
	}

	c.node.Remove(cid)

	atomic.AddInt64(&c.count, -1)
}
