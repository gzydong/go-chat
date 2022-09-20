package service

import (
	"context"
	"fmt"

	"go-chat/internal/repository/cache"
)

type ClientService struct {
	cache *cache.ClientStorage
}

func NewClientService(cache *cache.ClientStorage) *ClientService {
	return &ClientService{
		cache: cache,
	}
}

func (c *ClientService) Bind(ctx context.Context, channel string, clientId int64, uid int) {
	c.cache.Set(ctx, channel, fmt.Sprintf("%d", clientId), uid)
}

func (c *ClientService) UnBind(ctx context.Context, channel string, clientId int64) {
	c.cache.Del(ctx, channel, fmt.Sprintf("%d", clientId))
}
