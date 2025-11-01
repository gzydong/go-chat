package provider

import (
	"github.com/gzydong/go-chat/internal/repository/cache"
	"github.com/mojocn/base64Captcha"
)

func NewBase64Captcha(captcha *cache.CaptchaStorage) *base64Captcha.Captcha {
	return base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, captcha)
}
