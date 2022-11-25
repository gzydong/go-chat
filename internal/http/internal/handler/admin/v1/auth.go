package v1

import (
	"time"

	"github.com/mojocn/base64Captcha"
	"go-chat/api/pb/admin/v1"
	"go-chat/config"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
)

type Auth struct {
	config  *config.Config
	captcha *cache.CaptchaStorage
}

func NewAuth(config *config.Config, captcha *cache.CaptchaStorage) *Auth {
	return &Auth{config: config, captcha: captcha}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	params := &admin.AuthLoginRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.captcha.Verify(params.CaptchaVoucher, params.Captcha, true) {
		return ctx.InvalidParams("验证码填写不正确")
	}

	expiresAt := time.Now().Add(12 * time.Hour)

	// 生成登录凭证
	token := jwt.GenerateToken("admin", c.config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        "1",
		Issuer:    "im.admin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return ctx.Success(&admin.AuthLoginResponse{
		Auth: &admin.AccessToken{
			Type:        "Bearer",
			AccessToken: token,
			ExpiresIn:   60 * 60 * 12,
		},
	})
}

// Captcha 图形验证码
func (c *Auth) Captcha(ctx *ichat.Context) error {

	captcha := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, c.captcha)

	generate, base64, err := captcha.Generate()
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&admin.AuthCaptchaResponse{
		Voucher: generate,
		Captcha: base64,
	})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}
