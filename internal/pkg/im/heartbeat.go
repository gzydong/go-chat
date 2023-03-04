package im

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/timewheel"
)

const (
	heartbeatInterval = 30 // 心跳检测间隔时间
	heartbeatTimeout  = 75 // 心跳检测超时时间（超时时间是隔间检测时间的2.5倍以上）
)

var health *heartbeat

// 客户端心跳管理
type heartbeat struct {
	timeWheel *timewheel.SimpleTimeWheel
}

func init() {
	health = &heartbeat{}
	health.timeWheel = timewheel.NewSimpleTimeWheel(1*time.Second, 100, health.handle)
}

func (h *heartbeat) addClient(c *Client) {
	_ = h.timeWheel.Add(c, time.Duration(heartbeatTimeout)*time.Second)
}

func (h *heartbeat) delClient(c *Client) {
	h.timeWheel.Remove(c)
}

func (h *heartbeat) Start(ctx context.Context) error {

	go h.timeWheel.Start()

	for {
		select {
		case <-ctx.Done():
			h.timeWheel.Stop()
			return errors.New("heartbeat stop")
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (h *heartbeat) handle(timeWheel *timewheel.SimpleTimeWheel, value any) {
	c, ok := value.(*Client)
	if !ok {
		return
	}

	interval := int(time.Now().Unix() - c.lastTime)
	if interval > heartbeatTimeout {
		c.Close(2000, "心跳检测超时，连接自动关闭")
	} else if interval > heartbeatInterval {
		// 超过心跳间隔时间则主动推送一次消息
		_ = c.Write(&ClientOutContent{
			Content: jsonutil.Marshal(&Message{"heartbeat", "ping"}),
		})

		_ = timeWheel.Add(c, time.Duration(heartbeatInterval)*time.Second)
	}
}
