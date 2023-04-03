package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	redis *redis.Client
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{rds}
}

// Lock 获取 redis 锁
func (r *RedisLock) Lock(ctx context.Context, name string, expire int) bool {
	return r.redis.SetNX(ctx, r.name(name), 1, time.Duration(expire)*time.Second).Val()
}

// UnLock 释放 redis 锁
func (r *RedisLock) UnLock(ctx context.Context, name string) bool {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return false
	end`

	return r.redis.Eval(ctx, script, []string{r.name(name)}, 1).Err() == nil
}

func (r *RedisLock) name(name string) string {
	return fmt.Sprintf("im:lock:%s", name)
}
