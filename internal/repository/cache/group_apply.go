package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type GroupApplyStorage struct {
	redis *redis.Client
}

func NewGroupApplyStorage(rds *redis.Client) *GroupApplyStorage {
	return &GroupApplyStorage{rds}
}

func (g *GroupApplyStorage) Incr(ctx context.Context, uid int) {
	g.redis.Incr(ctx, g.name(uid))
}

func (g *GroupApplyStorage) Get(ctx context.Context, uid int) int {
	val, err := g.redis.Get(ctx, g.name(uid)).Int()
	if err != nil {
		return 0
	}

	return val
}

func (g *GroupApplyStorage) Del(ctx context.Context, uid int) {
	g.redis.Del(ctx, g.name(uid))
}

func (g *GroupApplyStorage) name(uid int) string {
	return fmt.Sprintf("im:group:apply:unread:uid_%d", uid)
}
