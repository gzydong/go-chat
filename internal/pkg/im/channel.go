package im

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/sourcegraph/conc/pool"
	"go-chat/internal/pkg/logger"
)

type IChannel interface {
	Name() string
	Count() int64
}

// Channel 渠道管理（多渠道划分，实现不同业务之间隔离）
type Channel struct {
	name    string // 渠道名称
	count   int64  // 客户端连接数
	node    *Node  // 客户端列表【客户端ID取余拆分，降低 map 长度】
	cmap    *cmap.ConcurrentMap[int64, *Client]
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
	return c.cmap.Get(cid)
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
	c.cmap.Set(client.cid, client)

	atomic.AddInt64(&c.count, 1)
}

// delClient 删除客户端
func (c *Channel) delClient(client *Client) {

	if !c.cmap.Has(client.cid) {
		return
	}

	c.cmap.Remove(client.cid)

	atomic.AddInt64(&c.count, -1)
}

// Start 渠道消费协程
func (c *Channel) Start(ctx context.Context) error {

	work := pool.New().WithMaxGoroutines(10)

	defer func() {
		log.Println(fmt.Errorf(fmt.Sprintf("loopPush 退出 %s", c.Name())))
		logger.Error(fmt.Sprintf("loopPush 退出 %s", c.Name()))
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("channel exit %s", c.name)
		case body, ok := <-c.outChan:
			if !ok {
				return fmt.Errorf(fmt.Sprintf("loopPush 退出 %s", c.Name()))
			}

			bodyContent := body
			content, _ := json.Marshal(bodyContent.GetMessage())

			work.Go(func() {
				if bodyContent.IsBroadcast() {
					c.cmap.IterCb(func(key int64, value *Client) {
						_ = value.Write(&ClientOutContent{Content: content})
					})
				} else {
					for _, cid := range bodyContent.receives {
						if client, ok := c.Client(cid); ok {
							_ = client.Write(&ClientOutContent{Content: content})
						}
					}
				}
			})
		}
	}
}
