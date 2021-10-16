package service

import (
	"context"
	"go-chat/app/cache"
)

type ClientService struct {
	cache *cache.WsClient
}

func NewClientService(cache *cache.WsClient) *ClientService {
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
