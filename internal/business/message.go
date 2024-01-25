package business

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
)

type PushMessage struct {
	Redis *redis.Client
}

func (m *PushMessage) Push(ctx context.Context, topic string, body *entity.SubscribeMessage) error {
	m.Redis.Publish(ctx, topic, jsonutil.Encode(body))
	return nil
}

func (m *PushMessage) MultiPush(ctx context.Context, topic string, items []*entity.SubscribeMessage) error {
	pipe := m.Redis.Pipeline()

	for _, body := range items {
		pipe.Publish(ctx, topic, jsonutil.Encode(body))
	}

	_, err := pipe.Exec(ctx)
	return err
}
