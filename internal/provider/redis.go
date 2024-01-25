package provider

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go-chat/config"
)

func NewRedisClient(conf *config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:        conf.Redis.Host,
		Password:    conf.Redis.Auth,
		DB:          conf.Redis.Database,
		ReadTimeout: 0,
	})

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		panic(fmt.Errorf("redis client error: %s", err))
	}

	return client
}
