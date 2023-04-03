package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/jsonutil"
)

const lastMessageCacheKey = "redis:hash:last-message"

type MessageStorage struct {
	redis *redis.Client
}

type LastCacheMessage struct {
	Content  string `json:"content"`
	Datetime string `json:"datetime"`
}

func NewMessageStorage(rds *redis.Client) *MessageStorage {
	return &MessageStorage{rds}
}

func (m *MessageStorage) Set(ctx context.Context, talkType int, sender int, receive int, message *LastCacheMessage) error {
	text := jsonutil.Encode(message)

	return m.redis.HSet(ctx, lastMessageCacheKey, m.name(talkType, sender, receive), text).Err()
}

func (m *MessageStorage) Get(ctx context.Context, talkType int, sender int, receive int) (*LastCacheMessage, error) {

	res, err := m.redis.HGet(ctx, lastMessageCacheKey, m.name(talkType, sender, receive)).Result()
	if err != nil {
		return nil, err
	}

	msg := &LastCacheMessage{}
	if err = jsonutil.Decode(res, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *MessageStorage) MGet(ctx context.Context, fields []string) ([]*LastCacheMessage, error) {

	res := m.redis.HMGet(ctx, lastMessageCacheKey, fields...)

	items := make([]*LastCacheMessage, 0)
	for _, item := range res.Val() {
		if val, ok := item.(string); ok {
			msg := &LastCacheMessage{}
			if err := jsonutil.Decode(val, msg); err != nil {
				return nil, err
			}

			items = append(items, msg)
		}
	}

	return items, nil
}

func (m *MessageStorage) name(talkType int, sender int, receive int) string {
	if talkType == 2 {
		sender = 0
	}

	if sender > receive {
		sender, receive = receive, sender
	}

	return fmt.Sprintf("%d_%d_%d", talkType, sender, receive)
}
