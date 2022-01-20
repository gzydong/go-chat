package process

import (
	"context"
	"go-chat/config"
	"go-chat/internal/cache"
	"log"
	"time"
)

type Server struct {
	conf   *config.Config
	server *cache.SidServer
}

func NewServer(conf *config.Config, server *cache.SidServer) *Server {
	return &Server{conf: conf, server: server}
}

func (s *Server) Handle(ctx context.Context) error {
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
