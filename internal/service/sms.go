package service

import (
	"context"
	"fmt"
	"time"

	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/cache"
)

var _ ISmsService = (*SmsService)(nil)

type ISmsService interface {
	Verify(ctx context.Context, channel string, mobile string, code string) bool
	Delete(ctx context.Context, channel string, mobile string)
	Send(ctx context.Context, channel string, mobile string) (string, error)
}

type SmsService struct {
	Storage *cache.SmsStorage
}

// Verify 验证短信验证码是否正确
func (s *SmsService) Verify(ctx context.Context, channel string, mobile string, code string) bool {
	return s.Storage.Verify(ctx, channel, mobile, code)
}

// Delete 删除短信验证码记录
func (s *SmsService) Delete(ctx context.Context, channel string, mobile string) {
	_ = s.Storage.Del(ctx, channel, mobile)
}

// Send 发送短信
func (s *SmsService) Send(ctx context.Context, channel string, mobile string) (string, error) {

	code := strutil.GenValidateCode(6)

	// 添加发送记录
	if err := s.Storage.Set(ctx, channel, mobile, code, 15*time.Minute); err != nil {
		return "", err
	}

	// TODO ... 请求第三方短信接口
	fmt.Println("正在发送短信验证码：", code)

	return code, nil
}
