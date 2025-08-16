package admin

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"go-chat/api/pb/admin/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/errorx"
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

	adminInfo, err := c.AdminRepo.FindByWhere(ctx.GetContext(), "email = ?", in.Username)
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

	if !adminInfo.VerifyPassword(string(password)) {
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
			c.ExpiresAt = jwt.NewNumericDate(expiresAt)
		},
	)

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&admin.AuthLoginResponse{
		Username: adminInfo.Username,
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

// Detail 退出登录接口
func (c *Auth) Detail(ctx *core.Context) error {
	info, err := c.AdminRepo.FindById(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return err
	}

	return ctx.Success(admin.AuthDetailResponse{
		Username:        info.Username,
		Email:           info.Email,
		Mobile:          info.Mobile,
		Address:         info.Address,
		TwoFactorEnable: "N",
	})
}

func (c *Auth) UpdatePassword(ctx *core.Context) error {
	var in admin.AuthUpdatePasswordRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	adminInfo, err := c.AdminRepo.FindById(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return err
	}

	oldPassword, err := c.Rsa.Decrypt(in.OldPassword)
	if err != nil {
		return ctx.Error(err)
	}

	newPassword, err := c.Rsa.Decrypt(in.NewPassword)
	if err != nil {
		return ctx.Error(err)
	}

	if !adminInfo.VerifyPassword(string(oldPassword)) {
		return errorx.New(400, "密码错误")
	}

	if string(oldPassword) == string(newPassword) {
		return errorx.New(400, "新密码不能与旧密码相同")
	}

	_, err = c.AdminRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
		"password": encrypt.HashPassword(string(newPassword)),
	}, "id = ?", ctx.AuthId())
	if err != nil {
		return err
	}

	return ctx.Success(nil)
}

func (c *Auth) UpdateDetail(ctx *core.Context) error {
	var in admin.AuthUpdateDetailRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := c.AdminRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
		"username": in.Username,
		"mobile":   in.Mobile,
		"address":  in.Address,
	}, "id = ?", ctx.AuthId())
	if err != nil {
		fmt.Println("======", err)
		return err
	}

	return ctx.Success(nil)
}
