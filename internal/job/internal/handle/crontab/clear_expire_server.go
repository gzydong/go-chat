package crontab

import (
	"context"

	"go-chat/internal/cache"
)

type ClearExpireServerHandle struct {
	server *cache.SidServer
}

func NewClearExpireServer(server *cache.SidServer) *ClearExpireServerHandle {
	return &ClearExpireServerHandle{server: server}
}

func (c *ClearExpireServerHandle) Handle(ctx context.Context) error {

	for _, sid := range c.server.All(ctx, 2) {
		_ = c.server.Del(context.Background(), sid)
		_ = c.server.SetExpireServer(context.Background(), sid)
	}

	return nil
}
