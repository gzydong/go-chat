package cache

import (
	"context"
	"fmt"
	"time"

	"go-chat/connect"
)

type AuthToken struct {
	Redis *connect.Redis
}

func (a *AuthToken) key(token string) string {
	return fmt.Sprintf("jwt:black-list:%s", token)
}

// SetBlackList 登录 token 加入黑名单
func (a *AuthToken) SetBlackList(ctx context.Context, token string, expiration int) error {
	ex := time.Duration(expiration) * time.Second

	return a.Redis.Client.Set(ctx, a.key(token), 1, ex).Err()
}

// DelBlackList 将 token 移出黑名单
func (a *AuthToken) DelBlackList(ctx context.Context, token string) error {
	return a.Redis.Client.Del(ctx, a.key(token)).Err()
}

// IsExistBlackList 判断 token 是否存在白名单
func (a AuthToken) IsExistBlackList(ctx context.Context, token string) error {
	return a.Redis.Client.Get(ctx, a.key(token)).Err()
}
