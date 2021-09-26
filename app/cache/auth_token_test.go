package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
)

func TestAuthToken_SetBlackList(t *testing.T) {
	client := testutil.TestRedisClient()
	s := &AuthToken{
		Redis: client,
	}

	ctx := context.Background()

	err := s.SetBlackList(ctx, "test-token", 10)
	assert.NoError(t, err)
	assert.NoError(t, s.IsExistBlackList(ctx, "test-token"))
}
