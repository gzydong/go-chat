package im

import (
	"context"
	"sync"
	"time"
)

const (
	heartbeatCheckInterval = 30 * time.Second // 心跳检测时间
	heartbeatIdleTime      = 70               // 心跳超时时间
	nodeLen                = 100              // 客户端拆分数量
)

var Heartbeat = &heartbeat{
	nodes: maps(nodeLen),
}

// 客户端心跳管理
type heartbeat struct {
	nodes []*sync.Map
}

func (h *heartbeat) node(cid int64) *sync.Map {
	return h.nodes[getMapIndex(cid, nodeLen)]
}

func (h *heartbeat) addClient(c *Client) {
	h.node(c.cid).Store(c.cid, c)
}

func (h *heartbeat) delClient(c *Client) {
	h.node(c.cid).Delete(c.cid)
}

func (h *heartbeat) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(heartbeatCheckInterval):
			t := time.Now().Unix()

			for _, node := range h.nodes {
				node.Range(func(key, value interface{}) bool {
					if c, ok := value.(*Client); ok && !c.IsClosed() {
						if int(t-c.lastTime) > heartbeatIdleTime {
							c.Close(2000, "心跳检测超时，连接自动关闭")
						}
					}

					return true
				})
			}
		}
	}
}
