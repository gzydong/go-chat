package cron

import (
	"context"

	"go-chat/internal/repository/cache"
)

type ClearExpireServer struct {
	storage *cache.SidStorage
}

func NewClearExpireServer(storage *cache.SidStorage) *ClearExpireServer {
	return &ClearExpireServer{storage: storage}
}

// Spec 配置定时任务规则
func (c *ClearExpireServer) Spec() string {
	return "* * * * *"
}

func (c *ClearExpireServer) Handle(ctx context.Context) error {

	for _, sid := range c.storage.All(ctx, 2) {
		_ = c.storage.Del(ctx, sid)
		_ = c.storage.SetExpireServer(ctx, sid)
	}

	return nil
}
