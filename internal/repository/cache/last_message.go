package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/pkg/jsonutil"
)

const lastMessageCacheKey = "rds:hash:last-message"

type LastMessage struct {
	rds *redis.Client
}

type LastCacheMessage struct {
	Content  string `json:"content"`
	Datetime string `json:"datetime"`
}

func NewLastMessage(rds *redis.Client) *LastMessage {
	return &LastMessage{rds}
}

func (cache *LastMessage) key(talkType int, sender int, receive int) string {
	if talkType == 2 {
		sender = 0
	}

	if sender > receive {
		sender, receive = receive, sender
	}

	return fmt.Sprintf("%d_%d_%d", talkType, sender, receive)
}

func (cache *LastMessage) Set(ctx context.Context, talkType int, sender int, receive int, message *LastCacheMessage) error {
	text := jsonutil.Encode(message)

	return cache.rds.HSet(ctx, lastMessageCacheKey, cache.key(talkType, sender, receive), text).Err()
}

func (cache *LastMessage) Get(ctx context.Context, talkType int, sender int, receive int) (*LastCacheMessage, error) {

	res, err := cache.rds.HGet(ctx, lastMessageCacheKey, cache.key(talkType, sender, receive)).Result()
	if err != nil {
		return nil, err
	}

	msg := &LastCacheMessage{}
	if err = jsonutil.Decode(res, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
