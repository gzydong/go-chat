package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/encrypt"
)

type SmsStorage struct {
	redis *redis.Client
}

func NewSmsStorage(redis *redis.Client) *SmsStorage {
	return &SmsStorage{redis}
}

func (s *SmsStorage) Set(ctx context.Context, channel string, mobile string, code string, exp time.Duration) error {
	_, err := s.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, s.failName(channel, mobile))
		pipe.Set(ctx, s.name(channel, mobile), code, exp)
		return nil
	})
	return err
}

func (s *SmsStorage) Get(ctx context.Context, channel string, mobile string) (string, error) {
	return s.redis.Get(ctx, s.name(channel, mobile)).Result()
}

func (s *SmsStorage) Del(ctx context.Context, channel string, mobile string) error {
	return s.redis.Del(ctx, s.name(channel, mobile)).Err()
}

func (s *SmsStorage) Verify(ctx context.Context, channel string, mobile string, code string) bool {

	value, err := s.Get(ctx, channel, mobile)
	if err != nil || len(value) == 0 {
		return false
	}

	if value == code {
		return true
	}

	// 3分钟内同一个手机号验证码错误次数超过5次，删除验证码
	num := s.redis.Incr(ctx, s.failName(channel, mobile)).Val()
	if num >= 5 {
		_, _ = s.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, s.name(channel, mobile))
			pipe.Del(ctx, s.failName(channel, mobile))
			return nil
		})
	} else if num == 1 {
		s.redis.Expire(ctx, s.failName(channel, mobile), 3*time.Minute)
	}

	return false
}

func (s *SmsStorage) name(channel string, mobile string) string {
	return fmt.Sprintf("im:auth:sms:%s:%s", channel, encrypt.Md5(mobile))
}

func (s *SmsStorage) failName(channel string, mobile string) string {
	return fmt.Sprintf("im:auth:sms_fail:%s:%s", channel, encrypt.Md5(mobile))
}
