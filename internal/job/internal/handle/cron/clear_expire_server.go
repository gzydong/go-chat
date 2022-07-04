package cron

import (
	"context"

	"go-chat/internal/repository/cache"
)

type ClearExpireServer struct {
	server *cache.SidServer
}

func NewClearExpireServer(server *cache.SidServer) *ClearExpireServer {
	return &ClearExpireServer{server: server}
}

// Spec 配置定时任务规则
func (c *ClearExpireServer) Spec() string {
	return "* * * * *"
}

func (c *ClearExpireServer) Handle(ctx context.Context) error {

	for _, sid := range c.server.All(ctx, 2) {
		_ = c.server.Del(ctx, sid)
		_ = c.server.SetExpireServer(ctx, sid)
	}

	return nil
}
