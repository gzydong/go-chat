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

func (s *Sequence) Name(id int, isUserId bool) string {
	if isUserId {
		return fmt.Sprintf("im:sequence:chat:uid:%d", id)
	}

	return fmt.Sprintf("im:sequence:chat:group:%d", id)
}

// Set 初始化发号器
func (s *Sequence) Set(ctx context.Context, id int, isUserId bool, value int64) error {
	return s.redis.SetEx(ctx, s.Name(id, isUserId), value, 12*time.Hour).Err()
}

// Get 获取消息时序ID
func (s *Sequence) Get(ctx context.Context, id int, isUserId bool) int64 {
	return s.redis.Incr(ctx, s.Name(id, isUserId)).Val()
}

// BatchGet 批量获取消息时序ID
func (s *Sequence) BatchGet(ctx context.Context, id int, isUserId bool, num int64) []int64 {

	value := s.redis.IncrBy(ctx, s.Name(id, isUserId), num).Val()

	items := make([]int64, 0, num)
	for i := num; i > 0; i-- {
		items = append(items, value-i+1)
	}

	return items
}
