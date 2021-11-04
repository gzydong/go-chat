package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/app/cache"
	"go-chat/testutil"
)

func testSmsService() *SmsService {
	codeCache := &cache.SmsCodeCache{Redis: testutil.TestRedisClient()}
	return &SmsService{smsCodeCache: codeCache}
}

func TestSmsService_SendSmsCode(t *testing.T) {
	s := testSmsService()
	err := s.SendSmsCode(context.Background(), "username", "1302012")
	assert.NoError(t, err)
}

func TestName(t *testing.T) {
	ids := []int{2, 3, 4, 5, 6, 7, 8}

	for _, id := range ids {
		fmt.Println(id)
	}
}
