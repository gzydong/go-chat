package service

import (
	"context"
	"fmt"
	"time"

	"go-chat/app/cache"
	"go-chat/app/helper"
)

type SmsService struct {
	SmsCodeCache *cache.SmsCodeCache
}

// CheckSmsCode 验证短信验证码是否正确
func (s *SmsService) CheckSmsCode(ctx context.Context, channel string, mobile string, code string) bool {
	value, err := s.SmsCodeCache.Get(ctx, channel, mobile)

	return err == nil && value == code
}

// DeleteSmsCode 删除短信验证码记录
func (s *SmsService) DeleteSmsCode(ctx context.Context, channel string, mobile string) {
	_ = s.SmsCodeCache.Del(ctx, channel, mobile)
}

// SendSmsCode 发送短信
func (s *SmsService) SendSmsCode(ctx context.Context, channel string, mobile string) error {
	// todo 需要做防止短信攻击处理

	code := helper.GenValidateCode(6)

	// 添加发送记录
	if err := s.SmsCodeCache.Set(ctx, channel, mobile, code, 10); err != nil {
		return err
	}

	// ... 请求第三方短信接口
	time.Sleep(2 * time.Second)
	fmt.Println("正在发送短信验证码：", code)

	return nil
}
