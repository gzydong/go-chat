package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenSessionStorage struct {
	rds *redis.Client
}

func NewTokenSessionStorage(rds *redis.Client) *TokenSessionStorage {
	return &TokenSessionStorage{rds}
}

func (s *TokenSessionStorage) key(token string) string {
	return fmt.Sprintf("jwt:black-list:%s", token)
}

// SetBlackList 登录 token 加入黑名单
func (s *TokenSessionStorage) SetBlackList(ctx context.Context, token string, expire time.Duration) error {
	return s.rds.Set(ctx, s.key(token), 1, expire).Err()
}

// DelBlackList 将 token 移出黑名单
func (s *TokenSessionStorage) DelBlackList(ctx context.Context, token string) error {
	return s.rds.Del(ctx, s.key(token)).Err()
}

// IsBlackList 判断 token 是否存在黑名单
func (s *TokenSessionStorage) IsBlackList(ctx context.Context, token string) bool {
	return s.rds.Get(ctx, s.key(token)).Val() != ""
}
