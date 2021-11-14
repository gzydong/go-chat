package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ServerRunIdKey = "server_ids"

	// ServerOverTime 运行检测超时时间（单位秒）
	ServerOverTime = 35
)

type ServerRunID struct {
	Redis *redis.Client
}

func NewServerRun(redis *redis.Client) *ServerRunID {
	return &ServerRunID{Redis: redis}
}

func (s *ServerRunID) SetServerID(ctx context.Context, server string, time int64) error {
	return s.Redis.HSet(ctx, ServerRunIdKey, server, time).Err()
}

// GetServerRunIdAll 获取指定状态的运行ID
// status 状态[1:运行中;2:已超时;3:全部]
func (s *ServerRunID) GetServerRunIdAll(ctx context.Context, status int) []string {
	result, err := s.Redis.HGetAll(ctx, ServerRunIdKey).Result()

	slice := make([]string, 0)

	t := time.Now().Unix()
	if err == nil {
		for key, val := range result {
			value, err := strconv.Atoi(val)

			if err != nil {
				continue
			}

			switch status {
			case 1:
				if t-int64(value) >= ServerOverTime {
					continue
				}
			case 2:
				if t-int64(value) < ServerOverTime {
					continue
				}
			}

			slice = append(slice, key)
		}
	}

	return slice
}
