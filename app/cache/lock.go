package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	Redis *redis.Client
}

// key 获取锁名
func (l *RedisLock) key(name string) string {
	return fmt.Sprintf("rds-lock:%s", name)
}

// Lock 获取 redis 分布式锁
func (l *RedisLock) Lock(ctx context.Context, name string, expire int) bool {
	return l.Redis.SetNX(ctx, l.key(name), 1, time.Duration(expire)*time.Second).Val()
}

// Release 释放 redis 分布式锁
func (l *RedisLock) Release(ctx context.Context, name string) bool {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return false
	end`

	return l.Redis.Eval(ctx, script, []string{l.key(name)}, 1).Err() == nil
}
