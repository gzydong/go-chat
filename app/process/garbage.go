package process

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClearGarbage struct {
	redis *redis.Client
}

func NewClearGarbage(redis *redis.Client) *ClearGarbage {
	return &ClearGarbage{redis: redis}
}

func (p *ClearGarbage) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Hour):

			items := p.redis.SMembers(ctx, "server_ids_expire").Val()

			for _, sid := range items {
				iter := p.redis.Scan(ctx, 0, fmt.Sprintf("ws:%s:*", sid), 100).Iterator()
				for iter.Next(ctx) {
					p.redis.Del(ctx, iter.Val())
				}

				p.redis.SRem(ctx, "server_ids_expire", sid)
			}
		}
	}
}
