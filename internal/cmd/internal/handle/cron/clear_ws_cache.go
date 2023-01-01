package cron

import (
	"context"
	"fmt"

	"go-chat/internal/repository/cache"
)

type ClearWsCache struct {
	storage *cache.ServerStorage
}

func NewClearWsCache(storage *cache.ServerStorage) *ClearWsCache {
	return &ClearWsCache{storage: storage}
}

// Spec 配置定时任务规则
// 每隔30分钟处理 websocket 缓存
func (c *ClearWsCache) Spec() string {
	return "*/30 * * * *"
}

func (c *ClearWsCache) Enable() bool {
	return true
}

func (c *ClearWsCache) Handle(ctx context.Context) error {

	for _, sid := range c.storage.GetExpireServerAll(ctx) {

		iter := c.storage.Redis().Scan(ctx, 0, fmt.Sprintf("ws:%s:*", sid), 100).Iterator()

		for iter.Next(ctx) {
			c.storage.Redis().Del(ctx, iter.Val())
		}

		_ = c.storage.DelExpireServer(ctx, sid)
	}

	return nil
}
