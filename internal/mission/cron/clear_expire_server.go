package cron

import (
	"context"

	"go-chat/internal/pkg/core/crontab"
	"go-chat/internal/repository/cache"
)

var _ crontab.ICrontab = (*ClearExpireServer)(nil)

type ClearExpireServer struct {
	Storage *cache.ServerStorage
}

func (c *ClearExpireServer) Name() string {
	return "expire.server.clear"
}

// Spec 配置定时任务规则
func (c *ClearExpireServer) Spec() string {
	return "*/10 * * * *"
}

func (c *ClearExpireServer) Enable() bool {
	return true
}

func (c *ClearExpireServer) Do(ctx context.Context) error {

	for _, sid := range c.Storage.All(ctx, 2) {
		_ = c.Storage.Del(ctx, sid)
		_ = c.Storage.SetExpireServer(ctx, sid)
	}

	return nil
}
