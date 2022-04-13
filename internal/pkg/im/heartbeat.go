package im

import (
	"context"
	"time"
)

const (
	heartbeatInterval = 30 // 心跳检测间隔时间
	heartbeatTimeout  = 75 // 心跳检测超时时间（超时时间是隔间检测时间的2.5倍以上）
)

var health *heartbeat

// 客户端心跳管理
type heartbeat struct {
	node *Node
}

func init() {
	health = &heartbeat{node: NewNode(10)}
}

func (h *heartbeat) addClient(c *Client) {
	h.node.add(c)
}

func (h *heartbeat) delClient(c *Client) {
	h.node.del(c)
}

func (h *heartbeat) Start(ctx context.Context) error {

	timer := time.NewTimer(heartbeatInterval * time.Second)

	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			t := time.Now().Unix()

			h.node.each(func(c *Client) {
				if int(t-c.lastTime) > heartbeatTimeout {
					c.Close(2000, "心跳检测超时，连接自动关闭")
				}
			})

			timer.Reset(heartbeatInterval * time.Second)
		}
	}
}
