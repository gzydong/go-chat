package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/app/cache"
	"go-chat/testutil"
)

func testSmsService() *SmsService {
	codeCache := &cache.SmsCodeCache{Redis: testutil.TestRedisClient()}
	return &SmsService{SmsCodeCache: codeCache}
}

func TestSmsService_SendSmsCode(t *testing.T) {
	s := testSmsService()
	err := s.SendSmsCode(context.Background(), "username", "1302012")
	assert.NoError(t, err)
}
