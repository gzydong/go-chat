package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	rds *redis.Client
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{rds}
}

// Lock 获取 redis 锁
func (lock *RedisLock) Lock(ctx context.Context, name string, expire int) bool {
	return lock.rds.SetNX(ctx, lock.name(name), 1, time.Duration(expire)*time.Second).Val()
}

// UnLock 释放 redis 锁
func (lock *RedisLock) UnLock(ctx context.Context, name string) bool {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return false
	end`

	return lock.rds.Eval(ctx, script, []string{lock.name(name)}, 1).Err() == nil
}

// 获取锁名
func (lock *RedisLock) name(name string) string {
	return fmt.Sprintf("redis:lock:%s", name)
}
