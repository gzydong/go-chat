package process

import (
	"context"

	"go-chat/internal/pkg/im"
)

// Heartbeat IM 客户端心跳检测管理
type Heartbeat struct {
}

func NewImHeartbeat() *Heartbeat {
	return &Heartbeat{}
}

// Handle 执行入口
func (s *Heartbeat) Handle(ctx context.Context) error {
	return im.Heartbeat.Run(ctx)
}
