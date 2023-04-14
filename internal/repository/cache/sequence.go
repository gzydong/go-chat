package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Sequence struct {
	redis *redis.Client
}

func NewSequence(redis *redis.Client) *Sequence {
	return &Sequence{redis: redis}
}

func (s *Sequence) Redis() *redis.Client {
	return s.redis
}

func (s *Sequence) Name(userId int, receiverId int) string {

	if userId == 0 {
		return fmt.Sprintf("im:sequence:chat:%d", receiverId)
	}

	if receiverId < userId {
		receiverId, userId = userId, receiverId
	}

	return fmt.Sprintf("im:sequence:chat:%d_%d", userId, receiverId)
}

// Set 初始化发号器
func (s *Sequence) Set(ctx context.Context, userId int, receiverId int, value int64) error {
	return s.redis.SetEx(ctx, s.Name(userId, receiverId), value, 12*time.Hour).Err()
}

// Get 获取消息时序ID
func (s *Sequence) Get(ctx context.Context, userId int, receiverId int) int64 {
	return s.redis.Incr(ctx, s.Name(userId, receiverId)).Val()
}

// BatchGet 批量获取消息时序ID
func (s *Sequence) BatchGet(ctx context.Context, userId int, receiverId int, num int64) []int64 {

	value := s.redis.IncrBy(ctx, s.Name(userId, receiverId), num).Val()

	items := make([]int64, 0, num)
	for i := num; i > 0; i-- {
		items = append(items, value-i+1)
	}

	return items
}
