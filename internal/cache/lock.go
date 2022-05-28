package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLock struct {
	rds *redis.Client
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{rds}
}

// key 获取锁名
func (lock *RedisLock) key(name string) string {
	return fmt.Sprintf("rds-lock:%s", name)
}

// Lock 获取 rds 锁
func (lock *RedisLock) Lock(ctx context.Context, name string, expire int) bool {
	return lock.rds.SetNX(ctx, lock.key(name), 1, time.Duration(expire)*time.Second).Val()
}

// UnLock 释放 rds 锁
func (lock *RedisLock) UnLock(ctx context.Context, name string) bool {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return false
	end`

	return lock.rds.Eval(ctx, script, []string{lock.key(name)}, 1).Err() == nil
}
