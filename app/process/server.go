package process

import (
	"context"
	"go-chat/app/cache"
	"go-chat/config"
	"log"
	"time"
)

type ServerRun struct {
	conf   *config.Config
	server *cache.ServerRunID
}

func NewServerRun(conf *config.Config, server *cache.ServerRunID) *ServerRun {
	return &ServerRun{conf, server}
}

func (s *ServerRun) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.server.SetServerID(ctx, s.conf.Server.ServerId, time.Now().Unix()); err != nil {
				log.Printf("SetServerID Error: %s\n", err)
			}
		}
	}
}
