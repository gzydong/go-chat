package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SmsCodeCache struct {
	rds *redis.Client
}

func NewSmsCodeCache(rds *redis.Client) *SmsCodeCache {
	return &SmsCodeCache{rds: rds}
}

func (c *SmsCodeCache) name(channel string, mobile string) string {
	return fmt.Sprintf("sms:%s:%s", channel, mobile)
}

func (c *SmsCodeCache) Set(ctx context.Context, channel string, mobile string, code string, exp time.Duration) error {
	// 发送新的短信，则清空失败次数
	c.rds.Del(ctx, fmt.Sprintf("sms:verify_fail:%s:%s", channel, mobile))
	return c.rds.Set(ctx, c.name(channel, mobile), code, exp).Err()
}

func (c *SmsCodeCache) Get(ctx context.Context, channel string, mobile string) (string, error) {
	return c.rds.Get(ctx, c.name(channel, mobile)).Result()
}

func (c *SmsCodeCache) Del(ctx context.Context, channel string, mobile string) error {
	return c.rds.Del(ctx, c.name(channel, mobile)).Err()
}

func (c *SmsCodeCache) IncrVerifyFail(ctx context.Context, channel string, mobile string, exp time.Duration) int64 {

	key := fmt.Sprintf("sms:verify_fail:%s:%s", channel, mobile)

	num := c.rds.Incr(ctx, key).Val()

	c.rds.Expire(ctx, key, exp)

	return num
}
