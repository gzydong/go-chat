package process

import (
	"context"
	"go-chat/app/pkg/im"
)

type Heartbeat struct {
}

func NewImHeartbeat() *Heartbeat {
	return &Heartbeat{}
}

// IM 客户端心跳检测管理
func (s *Heartbeat) Handle(ctx context.Context) error {

	im.Heartbeat.Run()

	return nil
}
