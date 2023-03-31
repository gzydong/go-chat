package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/encrypt"
)

type SmsStorage struct {
	redis *redis.Client
}

func NewSmsStorage(rds *redis.Client) *SmsStorage {
	return &SmsStorage{redis: rds}
}

func (c *SmsStorage) Set(ctx context.Context, channel string, mobile string, code string, exp time.Duration) error {
	c.redis.Del(ctx, c.failName(channel, mobile))
	return c.redis.Set(ctx, c.name(channel, mobile), code, exp).Err()
}

func (c *SmsStorage) Get(ctx context.Context, channel string, mobile string) (string, error) {
	return c.redis.Get(ctx, c.name(channel, mobile)).Result()
}

func (c *SmsStorage) Del(ctx context.Context, channel string, mobile string) error {
	return c.redis.Del(ctx, c.name(channel, mobile)).Err()
}

func (c *SmsStorage) Verify(ctx context.Context, channel string, mobile string, code string) bool {

	value, err := c.Get(ctx, channel, mobile)
	if err != nil {
		return false
	}

	if value == code {
		return true
	}

	// 3分钟内同一个手机号验证码错误次数超过5次，删除验证码
	num := c.redis.Incr(ctx, c.failName(channel, mobile)).Val()
	if num >= 5 {
		_ = c.Del(ctx, channel, mobile)
		c.redis.Del(ctx, c.failName(channel, mobile))
	} else if num == 1 {
		c.redis.Expire(ctx, c.failName(channel, mobile), 3*time.Minute)
	}

	return false
}

func (c *SmsStorage) name(channel string, mobile string) string {
	return fmt.Sprintf("sms:%s:%s", channel, encrypt.Md5(mobile))
}

func (c *SmsStorage) failName(channel string, mobile string) string {
	return fmt.Sprintf("sms:verify_fail:%s:%s", channel, encrypt.Md5(mobile))
}
