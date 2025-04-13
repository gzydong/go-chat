package v1

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"go-chat/api/pb/admin/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type Auth struct {
	Config          *config.Config
	AdminRepo       *repo.Admin
	JwtTokenStorage *cache.JwtTokenStorage
	ICaptcha        *base64Captcha.Captcha
	Rsa             rsautil.IRsa
}

// Login 登录接口
func (c *Auth) Login(ctx *core.Context) error {

	var in admin.AuthLoginRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.ICaptcha.Verify(in.CaptchaVoucher, in.Captcha, true) {
		return ctx.InvalidParams("验证码填写不正确")
	}

	adminInfo, err := c.AdminRepo.FindByWhere(ctx.GetContext(), "username = ? or email = ?", in.Username, in.Username)
	if err != nil {
		if utils.IsSqlNoRows(err) {
			return ctx.InvalidParams("账号不存在或密码填写错误!")
		}

		return ctx.Error(err)
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if !encrypt.VerifyPassword(adminInfo.Password, string(password)) {
		return ctx.InvalidParams("账号不存在或密码填写错误!")
	}

	if adminInfo.Status != model.AdminStatusNormal {
		return ctx.Error(entity.ErrAccountDisabled)
	}

	expiresAt := time.Now().Add(12 * time.Hour)

	token, err := jwtutil.NewTokenWithClaims(
		[]byte(c.Config.Jwt.Secret),
		entity.AdminClaims{
			AdminId: int32(adminInfo.Id),
		},
		func(c *jwt.RegisteredClaims) {
			c.Issuer = entity.JwtIssuerAdmin
		},
	)

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&admin.AuthLoginResponse{
		Auth: &admin.AccessToken{
			Type:        "Bearer",
			AccessToken: token,
			ExpiresIn:   int32(expiresAt.Unix() - time.Now().Unix()),
		},
	})
}

// Captcha 图形验证码
func (c *Auth) Captcha(ctx *core.Context) error {
	voucher, captcha, _, err := c.ICaptcha.Generate()
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&admin.AuthCaptchaResponse{
		Voucher: voucher,
		Captcha: captcha,
	})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *core.Context) error {
	token := middleware.GetAuthToken(ctx.Context)

	claims, err := jwtutil.ParseWithClaims[entity.AdminClaims]([]byte(c.Config.Jwt.Secret), token)
	if err == nil {
		if ex := claims.ExpiresAt.Unix() - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.GetContext(), token, time.Duration(ex)*time.Second)
		}
	}

	return ctx.Success(nil)
}
