package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
	"testing"
	"time"
)

func TestRedisLock_Lock(t *testing.T) {
	client := testutil.TestRedisClient()
	s := &RedisLock{
		Redis: client,
	}

	ctx := context.Background()

	res := s.Lock(ctx, "test-token", 20)
	assert.Equal(t, true, res)

	res2 := s.Lock(ctx, "test-token", 10)
	assert.Equal(t, false, res2)

	time.Sleep(21 * time.Second)

	res3 := s.Lock(ctx, "test-token", 10)
	assert.Equal(t, true, res3)
}
