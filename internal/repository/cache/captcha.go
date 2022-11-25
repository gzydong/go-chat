package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CaptchaStorage struct {
	redis *redis.Client
}

func NewCaptchaStorage(redis *redis.Client) *CaptchaStorage {
	return &CaptchaStorage{redis: redis}
}

func (c *CaptchaStorage) name(id string) string {
	return fmt.Sprintf("im:auth:captcha:%s", id)
}

func (c *CaptchaStorage) Set(id string, value string) error {
	return c.redis.SetEX(context.Background(), c.name(id), value, 3*time.Minute).Err()
}

func (c *CaptchaStorage) Get(id string, clear bool) string {

	value := c.redis.Get(context.Background(), c.name(id)).Val()
	if clear && len(value) > 0 {
		c.redis.Del(context.Background(), c.name(id))
	}

	return value
}

func (c *CaptchaStorage) Verify(id, answer string, clear bool) bool {

	value := c.redis.Get(context.Background(), c.name(id)).Val()
	if clear && len(value) > 0 {
		c.redis.Del(context.Background(), c.name(id))
	}

	return value == answer
}
