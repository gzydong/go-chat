package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisLock struct {
	Redis *redis.Client
}

func (l *RedisLock) key(name string) string {
	return fmt.Sprintf("lock:%s", name)
}

// Lock 获取 redis 分布式锁
func (l *RedisLock) Lock(ctx context.Context, name string, expire int) bool {
	err := l.Redis.Do(ctx, "set", l.key(name), 10, "ex", expire, "nx").Err()

	return err == nil
}

// Release 释放 redis 分布式锁
func (l *RedisLock) Release(ctx context.Context, name string) bool {
	return l.Redis.Del(ctx, l.key(name)).Val() > 0
}
