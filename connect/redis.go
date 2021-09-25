package connect

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go-chat/config"
)

type Redis struct {
	Prefix string
	Client *redis.Client
}

func RedisConnect(ctx context.Context, conf *config.Config) *Redis {

	// 建立连接
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Auth,
		DB:       conf.Redis.Database,
	})

	// 检测心跳
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("redis clint error: %s", err))
	}

	return &Redis{Client: client, Prefix: conf.Redis.Prefix}
}
