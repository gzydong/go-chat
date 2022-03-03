package process

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"go-chat/internal/cache"
)

type ClearGarbage struct {
	redis  *redis.Client
	lock   *cache.RedisLock
	server *cache.SidServer
}

// NewClearGarbage 清除 Websocket 相关过期垃圾数据
func NewClearGarbage(redis *redis.Client, lock *cache.RedisLock, server *cache.SidServer) *ClearGarbage {
	return &ClearGarbage{redis: redis, lock: lock, server: server}
}

// Handle 执行入口
func (s *ClearGarbage) Handle(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Minute)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, sid := range s.server.GetExpireServerAll(ctx) {
				iter := s.server.Redis().Scan(ctx, 0, fmt.Sprintf("ws:%s:*", sid), 100).Iterator()

				for iter.Next(ctx) {
					s.server.Redis().Del(ctx, iter.Val())
				}

				_ = s.server.DelExpireServer(ctx, sid)
			}
		}
	}
}
