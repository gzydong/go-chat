package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type SmsCodeCache struct {
	rds *redis.Client
}

func NewSmsCodeCache(rds *redis.Client) *SmsCodeCache {
	return &SmsCodeCache{rds: rds}
}

func (c *SmsCodeCache) key(channel string, mobile string) string {
	return fmt.Sprintf("sms:%s:%s", channel, mobile)
}

func (c *SmsCodeCache) Set(ctx context.Context, channel string, mobile string, code string, expire int) error {
	return c.rds.Set(ctx, c.key(channel, mobile), code, time.Duration(expire)*time.Second).Err()
}

func (c *SmsCodeCache) Get(ctx context.Context, channel string, mobile string) (string, error) {
	return c.rds.Get(ctx, c.key(channel, mobile)).Result()
}

func (c *SmsCodeCache) Del(ctx context.Context, channel string, mobile string) error {
	return c.rds.Del(ctx, c.key(channel, mobile)).Err()
}
