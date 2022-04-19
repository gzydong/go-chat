package cron

import (
	"context"
	"fmt"

	"go-chat/internal/cache"
)

type ClearWsCacheHandle struct {
	server *cache.SidServer
}

func NewClearWsCacheHandle(server *cache.SidServer) *ClearWsCacheHandle {
	return &ClearWsCacheHandle{server: server}
}

// Spec 配置定时任务规则
// 每隔30分钟处理 websocket 缓存
func (c *ClearWsCacheHandle) Spec() string {
	return "*/30 * * * *"
}

func (c *ClearWsCacheHandle) Handle(ctx context.Context) error {

	for _, sid := range c.server.GetExpireServerAll(ctx) {

		iter := c.server.Redis().Scan(ctx, 0, fmt.Sprintf("ws:%s:*", sid), 100).Iterator()

		for iter.Next(ctx) {
			c.server.Redis().Del(ctx, iter.Val())
		}

		_ = c.server.DelExpireServer(ctx, sid)
	}

	return nil
}
