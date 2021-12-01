package im

import (
	"sync"
	"time"
)

const (
	heartbeatCheckInterval = 30 * time.Second // 心跳检测时间
	heartbeatIdleTime      = 70               // 心跳超时时间
	nodeLen                = 100
)

var Heartbeat *heartbeat

type heartbeat struct {
	nodes []*sync.Map
}

func (h *heartbeat) node(cid int64) *sync.Map {
	return h.nodes[getMapIndex(cid, nodeLen)]
}

func (h *heartbeat) addClient(c *Client) {
	h.node(c.ClientId).Store(c.ClientId, c)
}

func (h *heartbeat) delClient(c *Client) {
	h.node(c.ClientId).Delete(c.ClientId)
}

func (h *heartbeat) run() {
	for {
		<-time.After(heartbeatCheckInterval)

		t := time.Now().Unix()

		for _, node := range h.nodes {
			node.Range(func(key, value interface{}) bool {
				if c, ok := value.(*Client); ok {
					if int(t-c.LastTime) > heartbeatIdleTime {
						c.Close(2000, "心跳检测超时，连接自动关闭")
					}
				}

				return true
			})
		}
	}
}

func init() {
	Heartbeat = &heartbeat{
		nodes: maps(nodeLen),
	}

	go Heartbeat.run()
}
