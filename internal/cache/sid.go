package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// ServerKey 正在的运行服务
	ServerKey = "server_ids"

	// ServerKeyExpire 过期的运行服务
	ServerKeyExpire = "server_ids_expire"

	// ServerOverTime 运行检测超时时间（单位秒）
	ServerOverTime = 50
)

type SidServer struct {
	rds *redis.Client
}

func NewSid(rds *redis.Client) *SidServer {
	return &SidServer{rds: rds}
}

// Set 更新服务心跳时间
func (s *SidServer) Set(ctx context.Context, server string, time int64) error {

	_ = s.DelExpireServer(ctx, server)

	return s.rds.HSet(ctx, ServerKey, server, time).Err()
}

// Del 删除指定 SidServer
func (s *SidServer) Del(ctx context.Context, server string) error {
	return s.rds.HDel(ctx, ServerKey, server).Err()
}

// All 获取指定状态的运行 SidServer
// status 状态[1:运行中;2:已超时;3:全部]
func (s *SidServer) All(ctx context.Context, status int) []string {
	result, err := s.rds.HGetAll(ctx, ServerKey).Result()

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

func (s *SidServer) SetExpireServer(ctx context.Context, server string) error {
	return s.rds.SAdd(ctx, ServerKeyExpire, server).Err()
}

func (s *SidServer) DelExpireServer(ctx context.Context, server string) error {
	return s.rds.SRem(ctx, ServerKeyExpire, server).Err()
}

func (s *SidServer) GetExpireServerAll(ctx context.Context) []string {
	return s.rds.SMembers(ctx, ServerKeyExpire).Val()
}

func (s *SidServer) Redis() *redis.Client {
	return s.rds
}
