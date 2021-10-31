package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type UnreadTalkCache struct {
	rds *redis.Client
}

// Increment 消息未读数自增
// @params sender  发送者ID
// @params receive 接收者ID
func (c *UnreadTalkCache) Increment(ctx context.Context, sender int, receive int) {

}
