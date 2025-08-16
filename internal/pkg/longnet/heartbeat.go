package longnet

import (
	"log/slog"
	"time"

	"go-chat/internal/pkg/longnet/timewheel"
)

var _ IHeartbeat = (*Heartbeat)(nil)

type Heartbeat struct {
	tw *timewheel.TimingWheel
}

func NewHeartbeat(interval int32, fn func(taskId int64)) *Heartbeat {
	tw := timewheel.NewTimingWheel(time.Second, int(float32(interval)*1.5), 10)
	tw.SetCallback(fn)

	return &Heartbeat{
		tw: tw,
	}
}

func (h *Heartbeat) Insert(connId int64, d time.Duration) {
	err := h.tw.AddTask(connId, d, nil)
	if err != nil {
		slog.Error("add task error", "error", err)
	}
}

func (h *Heartbeat) Cancel(connId int64) {
	h.tw.Cancel(connId)
}
