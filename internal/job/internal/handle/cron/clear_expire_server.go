package cron

import (
	"context"

	"go-chat/internal/cache"
)

type ClearExpireServerHandle struct {
	server *cache.SidServer
}

func NewClearExpireServer(server *cache.SidServer) *ClearExpireServerHandle {
	return &ClearExpireServerHandle{server: server}
}

// Spec 配置定时任务规则
func (c *ClearExpireServerHandle) Spec() string {
	return "* * * * *"
}

func (c *ClearExpireServerHandle) Handle(ctx context.Context) error {

	for _, sid := range c.server.All(ctx, 2) {
		_ = c.server.Del(ctx, sid)
		_ = c.server.SetExpireServer(ctx, sid)
	}

	return nil
}
