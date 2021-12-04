package im

import (
	"context"
	"time"
)

const (
	heartbeatCheckInterval = 30 * time.Second // 心跳检测时间
	heartbeatIdleTime      = 70               // 心跳超时时间
)

var Heartbeat = &heartbeat{
	node: NewNode(100),
}

// 客户端心跳管理
type heartbeat struct {
	node *Node
}

func (h *heartbeat) addClient(c *Client) {
	h.node.add(c)
}

func (h *heartbeat) delClient(c *Client) {
	h.node.del(c)
}

func (h *heartbeat) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(heartbeatCheckInterval):
			t := time.Now().Unix()

			h.node.each(func(c *Client) {
				if int(t-c.lastTime) > heartbeatIdleTime {
					c.Close(2000, "心跳检测超时，连接自动关闭")
				}
			})
		}
	}
}
