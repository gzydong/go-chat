package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-chat/config"
)

var Rdb *redis.Client

// 初始化 Redis 连接
func NewRedis() {
	conf := config.GlobalConfig.Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Auth,
		DB:       conf.Database,
	})

	_, err := Rdb.Ping().Result()
	if err != nil {
		panic("Redis connection failed！")
	}
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() {
	_ = Rdb.Close()
}
