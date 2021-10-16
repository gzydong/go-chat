package service

import (
	"context"
	"go-chat/app/cache"
)

type ClientService struct {
	WsClient *cache.WsClient
}

// Bind ...
func (c *ClientService) Bind(ctx context.Context, channel string, clientId string, id int) {
	c.WsClient.Set(ctx, channel, clientId, id)
}

// UnBind ...
func (c *ClientService) UnBind(ctx context.Context, channel string, clientId string) {
	c.WsClient.Del(ctx, channel, clientId)
}
