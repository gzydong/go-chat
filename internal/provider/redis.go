package provider

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"go-chat/config"
)

func NewRedisClient(ctx context.Context, conf *config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:    conf.Redis.Auth,
		DB:          conf.Redis.Database,
		ReadTimeout: -1,
	})

	// 检测心跳
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("redis client error: %s", err))
	}

	return client
}
