package process

import (
	"context"
	"go-chat/internal/pkg/im"
)

type Heartbeat struct {
}

func NewImHeartbeat() *Heartbeat {
	return &Heartbeat{}
}

// IM 客户端心跳检测管理
func (s *Heartbeat) Handle(ctx context.Context) error {
	return im.Heartbeat.Run(ctx)
}
