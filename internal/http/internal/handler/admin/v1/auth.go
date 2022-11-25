package v1

import (
	"github.com/mojocn/base64Captcha"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/cache"
)

type Auth struct {
	captcha *cache.CaptchaStorage
}

func NewAuth(captcha *cache.CaptchaStorage) *Auth {
	return &Auth{captcha: captcha}
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

	return ctx.Success(&admin.AuthLoginResponse{})
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
