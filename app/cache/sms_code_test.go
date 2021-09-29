package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
	"testing"
)

func TestSmsCodeCache_Set(t *testing.T) {
	client := testutil.TestRedisClient()
	ctx := context.Background()

	s := &SmsCodeCache{
		Redis: client,
	}

	assert.NoError(t, s.Set(ctx, "register", "13888888888", "589641", 10))

	code, err := s.Get(ctx, "register", "13888888888")
	assert.NoError(t, err)
	assert.Equal(t, "589641", code)
}

func TestSmsCodeCache_Del(t *testing.T) {
	client := testutil.TestRedisClient()
	ctx := context.Background()

	s := &SmsCodeCache{
		Redis: client,
	}

	assert.NoError(t, s.Set(ctx, "register", "13888888888", "589641", 10))
	assert.NoError(t, s.Del(ctx, "register", "13888888888"))

	code, err := s.Get(ctx, "register", "13888888888")
	assert.Equal(t, "", code)
	assert.Error(t, err)
}
