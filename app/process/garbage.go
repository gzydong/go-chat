package process

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/cache"
	"time"
)

type ClearGarbage struct {
	redis  *redis.Client
	lock   *cache.RedisLock
	server *cache.SidServer
}

// 清除 Websocket 相关过期垃圾数据
func NewClearGarbage(redis *redis.Client, lock *cache.RedisLock) *ClearGarbage {
	return &ClearGarbage{redis: redis, lock: lock}
}

func (s *ClearGarbage) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Hour):
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
