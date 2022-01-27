package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Session struct {
	rds *redis.Client
}

func NewSession(rds *redis.Client) *Session {
	return &Session{rds}
}

func (a *Session) key(token string) string {
	return fmt.Sprintf("jwt:black-list:%s", token)
}

// SetBlackList 登录 token 加入黑名单
func (a *Session) SetBlackList(ctx context.Context, token string, expire int) error {
	ex := time.Duration(expire) * time.Second

	return a.rds.Set(ctx, a.key(token), time.Now().Add(time.Minute*3).Unix(), ex).Err()
}

// DelBlackList 将 token 移出黑名单
func (a *Session) DelBlackList(ctx context.Context, token string) error {
	return a.rds.Del(ctx, a.key(token)).Err()
}

// IsBlackList 判断 token 是否存在黑名单
func (a *Session) IsBlackList(ctx context.Context, token string) bool {
	val := a.rds.Get(ctx, a.key(token)).Val()

	if val == "" {
		return false
	}

	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false
	}

	// 判断是否在缓冲区时间内
	return res <= time.Now().Unix()
}
