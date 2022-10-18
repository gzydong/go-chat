package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Sequence struct {
	redis *redis.Client
}

func NewSequence(redis *redis.Client) *Sequence {
	return &Sequence{redis: redis}
}

// Seq 获取消息时序ID
func (s *Sequence) Seq(ctx context.Context, userId int, receiverId int) int64 {

	if receiverId < userId {
		receiverId, userId = userId, receiverId
	}

	return s.redis.Incr(ctx, fmt.Sprintf("im:sequence:%d_%d", userId, receiverId)).Val()
}
