package process

import (
	"context"
	"log"
	"time"

	"go-chat/config"
	"go-chat/internal/cache"
)

type Health struct {
	conf   *config.Config
	server *cache.SidServer
}

func NewHealthCheck(conf *config.Config, server *cache.SidServer) *Health {
	return &Health{conf: conf, server: server}
}

func (s *Health) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.server.SetServer(ctx, s.conf.ServerId(), time.Now().Unix()); err != nil {
				log.Printf("SetServer Error: %s\n", err)
				continue
			}

			for _, sid := range s.server.GetServerAll(ctx, 2) {
				_ = s.server.DelServer(ctx, sid)
				_ = s.server.SetExpireServer(ctx, sid)
			}
		}
	}
}
