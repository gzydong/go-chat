package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/connect"
	"go-chat/testutil"
)

func TestServerRunID_SetServerID(t *testing.T) {
	client := testutil.TestRedisClient()
	s := NewServerRun(&connect.Redis{Client: client})
	ctx := context.Background()
	assert.NoError(t, s.SetServerID(ctx, "jinxing.liu", 1))
	data, err := client.HGet(ctx, ServerRunIdKey, "jinxing.liu").Result()
	assert.NoError(t, err)
	assert.Equal(t, "1", data)
	fmt.Println(data, err)
}
