package process

import (
	"context"
	"log"
	"time"

	"go-chat/config"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

type HealthSubscribe struct {
	config  *config.Config
	storage *cache.ServerStorage
}

func NewHealthSubscribe(config *config.Config, storage *cache.ServerStorage) *HealthSubscribe {
	return &HealthSubscribe{config: config, storage: storage}
}

func (s *HealthSubscribe) Setup(ctx context.Context) error {

	log.Println("Start HealthSubscribe")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.storage.Set(ctx, s.config.ServerId(), time.Now().Unix()); err != nil {
				logger.Errorf("Websocket HealthSubscribe Report Err: %s \n", err.Error())
			}
		}
	}
}
