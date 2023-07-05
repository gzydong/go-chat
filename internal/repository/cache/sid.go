package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// ServerKey 正在的运行服务
	ServerKey = "server_ids"

	// ServerKeyExpire 过期的运行服务
	ServerKeyExpire = "server_ids_expire"

	// ServerOverTime 运行检测超时时间（单位秒）
	ServerOverTime = 50
)

type ServerStorage struct {
	redis *redis.Client
}

func NewSidStorage(rds *redis.Client) *ServerStorage {
	return &ServerStorage{redis: rds}
}

// Set 更新服务心跳时间
func (s *ServerStorage) Set(ctx context.Context, server string, time int64) error {
	_, err := s.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, ServerKeyExpire, server)
		pipe.HSet(ctx, ServerKey, server, time)
		return nil
	})
	return err
}

// Del 删除指定 ServerStorage
func (s *ServerStorage) Del(ctx context.Context, server string) error {
	return s.redis.HDel(ctx, ServerKey, server).Err()
}

// All 获取指定状态的运行 ServerStorage
// status 状态[1:运行中;2:已超时;3:全部]
func (s *ServerStorage) All(ctx context.Context, status int) []string {

	var (
		unix  = time.Now().Unix()
		slice = make([]string, 0)
	)

	all, err := s.redis.HGetAll(ctx, ServerKey).Result()
	if err != nil {
		return slice
	}

	for key, val := range all {
		value, err := strconv.Atoi(val)
		if err != nil {
			continue
		}

		switch status {
		case 1:
			if unix-int64(value) >= ServerOverTime {
				continue
			}
		case 2:
			if unix-int64(value) < ServerOverTime {
				continue
			}
		}

		slice = append(slice, key)
	}

	return slice
}

func (s *ServerStorage) SetExpireServer(ctx context.Context, server string) error {
	return s.redis.SAdd(ctx, ServerKeyExpire, server).Err()
}

func (s *ServerStorage) DelExpireServer(ctx context.Context, server string) error {
	return s.redis.SRem(ctx, ServerKeyExpire, server).Err()
}

func (s *ServerStorage) GetExpireServerAll(ctx context.Context) []string {
	return s.redis.SMembers(ctx, ServerKeyExpire).Val()
}

func (s *ServerStorage) Redis() *redis.Client {
	return s.redis
}
