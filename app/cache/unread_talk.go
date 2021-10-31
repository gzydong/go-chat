package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type UnreadTalkCache struct {
	rds *redis.Client
}

func NewUnreadTalkCache(rds *redis.Client) *UnreadTalkCache {
	return &UnreadTalkCache{rds}
}

// Increment 消息未读数自增
// @params sender  发送者ID
// @params receive 接收者ID
func (c *UnreadTalkCache) Increment(ctx context.Context, sender int, receive int) {

}
