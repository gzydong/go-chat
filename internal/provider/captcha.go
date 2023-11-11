package provider

import (
	"github.com/mojocn/base64Captcha"
	"go-chat/internal/repository/cache"
)

func NewBase64Captcha(captcha *cache.CaptchaStorage) *base64Captcha.Captcha {
	return base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, captcha)
}
