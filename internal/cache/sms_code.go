package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type SmsCodeCache struct {
	Redis *redis.Client
}

func (c *SmsCodeCache) key(channel string, mobile string) string {
	return fmt.Sprintf("sms:%s:%s", channel, mobile)
}

func (c *SmsCodeCache) Set(ctx context.Context, channel string, mobile string, code string, expire int) error {
	return c.Redis.Set(ctx, c.key(channel, mobile), code, time.Duration(expire)*time.Second).Err()
}

func (c *SmsCodeCache) Get(ctx context.Context, channel string, mobile string) (string, error) {
	return c.Redis.Get(ctx, c.key(channel, mobile)).Result()
}

func (c *SmsCodeCache) Del(ctx context.Context, channel string, mobile string) error {
	return c.Redis.Del(ctx, c.key(channel, mobile)).Err()
}
