package service

import (
	"context"
	"fmt"

	"go-chat/internal/cache"
	"go-chat/internal/pkg/strutil"
)

type SmsService struct {
	smsCodeCache *cache.SmsCodeCache
}

func NewSmsService(codeCache *cache.SmsCodeCache) *SmsService {
	return &SmsService{smsCodeCache: codeCache}
}

// CheckSmsCode 验证短信验证码是否正确
func (s *SmsService) CheckSmsCode(ctx context.Context, channel string, mobile string, code string) bool {
	value, err := s.smsCodeCache.Get(ctx, channel, mobile)

	return err == nil && value == code
}

// DeleteSmsCode 删除短信验证码记录
func (s *SmsService) DeleteSmsCode(ctx context.Context, channel string, mobile string) {
	_ = s.smsCodeCache.Del(ctx, channel, mobile)
}

// SendSmsCode 发送短信
func (s *SmsService) SendSmsCode(ctx context.Context, channel string, mobile string) (string, error) {
	// todo 需要做防止短信攻击处理

	code := strutil.GenValidateCode(6)

	// 添加发送记录
	if err := s.smsCodeCache.Set(ctx, channel, mobile, code, 60*15); err != nil {
		return "", err
	}

	// ... 请求第三方短信接口
	fmt.Println("正在发送短信验证码：", code)

	return code, nil
}
