package service

import (
	"context"

	"go-chat/internal/cache"
)

type ClientService struct {
	cache *cache.WsClientSession
}

func NewClientService(cache *cache.WsClientSession) *ClientService {
	return &ClientService{
		cache: cache,
	}
}

// Bind ...
func (c *ClientService) Bind(ctx context.Context, channel string, clientId string, id int) {
	c.cache.Set(ctx, channel, clientId, id)
}

// UnBind ...
func (c *ClientService) UnBind(ctx context.Context, channel string, clientId string) {
	c.cache.Del(ctx, channel, clientId)
}
