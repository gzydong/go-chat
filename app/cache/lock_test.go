package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
	"testing"
)

func TestRedisLock_Lock(t *testing.T) {
	client := testutil.TestRedisClient()
	s := &RedisLock{
		Redis: client,
	}

	ctx := context.Background()

	res := s.Lock(ctx, "test-token", 5)
	assert.Equal(t, true, res)

	res2 := s.Lock(ctx, "test-token", 5)
	assert.Equal(t, false, res2)
}
