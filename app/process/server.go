package process

import (
	"context"
	"go-chat/app/cache"
	"go-chat/config"
	"log"
	"time"
)

type Server struct {
	conf   *config.Config
	server *cache.Server
}

func NewServerRun(conf *config.Config, server *cache.Server) *Server {
	return &Server{conf: conf, server: server}
}

func (s *Server) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.server.SetServer(ctx, s.conf.GetSid(), time.Now().Unix()); err != nil {
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
