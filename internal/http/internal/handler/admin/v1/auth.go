package v1

import (
	"fmt"
	"time"

	"github.com/mojocn/base64Captcha"
	"go-chat/api/pb/admin/v1"
	"go-chat/config"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
)

type Auth struct {
	config  *config.Config
	captcha *cache.CaptchaStorage
	test    *repo.Test
}

func NewAuth(config *config.Config, captcha *cache.CaptchaStorage, test *repo.Test) *Auth {
	return &Auth{config: config, captcha: captcha, test: test}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	return ctx.ErrorBusiness("0000")

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

	// all, err := c.test.FindAll(ctx.Ctx(), func(db *gorm.DB) {
	// 	db.Order("id desc").Limit(3)
	// })
	// if err != nil {
	// 	return ctx.Error(err.Error())
	// }
	//
	// for _, m := range all {
	// 	fmt.Println(m.Id)
	// 	fmt.Printf("%T \n", m)
	// }
	//
	// data, err := c.test.FindById(ctx.Ctx(), 10)
	// if err != nil {
	// 	return ctx.Error(err.Error())
	// }
	//
	// fmt.Println(jsonutil.Encode(data))
	//
	// data, err = c.test.FindByWhere(ctx.Ctx(), "user_id = ? and class_id = ?", 4135, 3236)
	// if err != nil {
	// 	return ctx.Error(err.Error())
	// }
	//
	// fmt.Println(jsonutil.Encode(data))

	// count, err := c.test.UpdateWhere(ctx.Ctx(), map[string]interface{}{
	// 	"is_asterisk": 0,
	// }, "1 = 1")
	//
	// fmt.Println(count, err)

	isTrue, err := c.test.QueryExist(ctx.Ctx(), "user_id = ?", 2054)
	if err != nil {
		return err
	}

	fmt.Println(isTrue)

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *ichat.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}
