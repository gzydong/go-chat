package im

import (
	"context"
	"time"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/worker"
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

			h.check()

			timer.Reset(heartbeatInterval * time.Second)
		}
	}
}

func (h *heartbeat) check() {

	work := worker.NewTask(10)

	for _, val := range h.node.nodes {
		node := val

		work.Do(func() {

			ctime := time.Now().Unix()

			node.Range(func(key, value interface{}) bool {
				c := value.(*Client)

				interval := int(ctime - c.lastTime)
				if interval > heartbeatTimeout {
					c.Close(2000, "心跳检测超时，连接自动关闭")
				} else if interval > heartbeatInterval {
					// 超过心跳间隔时间则主动推送一次消息
					_ = c.Write(&ClientOutContent{
						Content: jsonutil.EncodeToBt(&Message{"heartbeat", "ping"}),
					})
				}

				return true
			})
		})
	}

	work.Wait()
}
