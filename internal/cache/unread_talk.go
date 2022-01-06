package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type UnreadTalkCache struct {
	rds *redis.Client
}

func NewUnreadTalkCache(rds *redis.Client) *UnreadTalkCache {
	return &UnreadTalkCache{rds}
}

func (c *UnreadTalkCache) key(sender, receive int) string {
	return fmt.Sprintf("%d_%d", sender, receive)
}

// Increment 消息未读数自增
// @params sender  发送者ID
// @params receive 接收者ID
func (c *UnreadTalkCache) Increment(ctx context.Context, sender, receive int) {
	c.rds.HIncrBy(ctx, "talk:unread:msg", c.key(sender, receive), 1)
}

// Get 获取消息未读数
// @params sender  发送者ID
// @params receive 接收者ID
func (c *UnreadTalkCache) Get(ctx context.Context, sender, receive int) int {
	val, _ := c.rds.HGet(ctx, "talk:unread:msg", c.key(sender, receive)).Int()

	return val
}

func (c *UnreadTalkCache) Reset(ctx context.Context, sender, receive int) {
	c.rds.HSet(ctx, "talk:unread:msg", c.key(sender, receive), 0)
}
