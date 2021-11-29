package process

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-chat/app/cache"
	"go-chat/config"
	"log"
	"time"
)

type ServerRun struct {
	conf   *config.Config
	server *cache.ServerRunID
	redis  *redis.Client
}

func NewServerRun(conf *config.Config, server *cache.ServerRunID, redis *redis.Client) *ServerRun {
	return &ServerRun{conf: conf, server: server, redis: redis}
}

func (s *ServerRun) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			if err := s.server.SetServerID(ctx, s.conf.GetSid(), time.Now().Unix()); err != nil {
				log.Printf("SetServerID Error: %s\n", err)
				continue
			}

			for _, sid := range s.server.GetServerRunIdAll(ctx, 2) {
				// iter := s.redis.Scan(ctx, 0, fmt.Sprintf("ws:%s:*", sid), 10).Iterator()
				// for iter.Next(ctx) {
				// 	s.redis.Del(ctx, iter.Val())
				// }

				_ = s.server.Del(ctx, sid)

				s.redis.SAdd(ctx, "server_ids_expire", sid)
			}
		}
	}
}
