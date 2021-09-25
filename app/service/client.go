package service

import (
	"context"

	"go-chat/app/cache"
)

type ClientService struct {
	WsClient *cache.WsClient
}

// Bind ...
func (c *ClientService) Bind(ctx context.Context, channel string, uuid string, id int) {
	c.WsClient.Set(ctx, channel, uuid, id)
}

// UnBind ...
func (c *ClientService) UnBind(ctx context.Context, channel string, uuid string) {
	c.WsClient.Del(ctx, channel, uuid)
}
