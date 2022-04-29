package process

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go-chat/config"
	"go-chat/internal/cache"
)

type Health struct {
	conf   *config.Config
	server *cache.SidServer
}

func NewHealth(conf *config.Config, server *cache.SidServer) *Health {
	return &Health{conf: conf, server: server}
}

func (s *Health) Setup(ctx context.Context) error {
	for {
		select {

		case <-ctx.Done():
			return nil

		// 每隔10秒上报心跳
		case <-time.After(10 * time.Second):
			if err := s.server.Set(ctx, s.conf.ServerId(), time.Now().Unix()); err != nil {
				logrus.Errorf("Websocket Health Report Err: %s \n", err.Error())
			}
		}
	}
}
