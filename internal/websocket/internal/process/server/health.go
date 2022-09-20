package server

import (
	"context"
	"log"
	"time"

	"go-chat/config"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

type HealthSubscribe struct {
	conf   *config.Config
	server *cache.SidServer
}

func NewHealthSubscribe(conf *config.Config, server *cache.SidServer) *HealthSubscribe {
	return &HealthSubscribe{conf: conf, server: server}
}

func (s *HealthSubscribe) Setup(ctx context.Context) error {

	log.Println("Start HealthSubscribe")

	for {
		select {

		case <-ctx.Done():
			return nil

		// 每隔10秒上报心跳
		case <-time.After(10 * time.Second):
			if err := s.server.Set(ctx, s.conf.ServerId(), time.Now().Unix()); err != nil {
				logger.Errorf("Websocket HealthSubscribe Report Err: %s \n", err.Error())
			}
		}
	}
}
