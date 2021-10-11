package service

import (
	"context"
	"log"
	"time"

	"go-chat/app/cache"
	"go-chat/config"
)

type SocketService struct {
	Conf        *config.Config
	ServerRunID *cache.ServerRunID
}

func (s *SocketService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.ServerRunID.SetServerID(ctx, s.Conf.Server.ServerId, time.Now().Unix()); err != nil {
				log.Printf("SetServerID Error: %s\n", err)
			}
		}
	}
}
