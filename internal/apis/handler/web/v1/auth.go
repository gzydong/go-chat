package v1

import (
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"go-chat/api/pb/queue/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Auth struct {
	Config              *config.Config
	Redis               *redis.Client
	JwtTokenStorage     *cache.JwtTokenStorage
	RedisLock           *cache.RedisLock
	RobotRepo           *repo.Robot
	SmsService          service.ISmsService
	UserService         service.IUserService
	ArticleClassService service.IArticleClassService
}

// Login 登录接口
func (c *Auth) Login(ctx *core.Context) error {
	in := &web.AuthLoginRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.UserService.Login(in.Mobile, in.Password)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	data := jsonutil.Marshal(queue.UserLoginRequest{
		UserId:   int32(user.Id),
		IpAddr:   ctx.Context.ClientIP(),
		Platform: in.Platform,
		Agent:    ctx.Context.GetHeader("user-agent"),
		LoginAt:  time.Now().Format(time.DateTime),
	})

	if err := c.Redis.Publish(ctx.Ctx(), entity.LoginTopic, data).Err(); err != nil {
		logger.ErrorWithFields(
			"投递登录消息异常", err,
			queue.UserLoginRequest{
				UserId:   int32(user.Id),
				IpAddr:   ctx.Context.ClientIP(),
				Platform: in.Platform,
				Agent:    ctx.Context.GetHeader("user-agent"),
				LoginAt:  time.Now().Format(time.DateTime),
			},
		)
	}

	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Bearer",
		AccessToken: c.token(user.Id),
		ExpiresIn:   int32(c.Config.Jwt.ExpiresTime),
	})
}

// Register 注册接口
func (c *Auth) Register(ctx *core.Context) error {
	in := &web.AuthRegisterRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.Ctx(), entity.SmsRegisterChannel, in.Mobile, in.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.UserService.Register(ctx.Ctx(), &service.UserRegisterOpt{
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Password: in.Password,
		Platform: in.Platform,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	c.SmsService.Delete(ctx.Ctx(), entity.SmsRegisterChannel, in.Mobile)

	return ctx.Success(&web.AuthRegisterResponse{})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *core.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *core.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(&web.AuthRefreshResponse{
		Type:        "Bearer",
		AccessToken: c.token(ctx.UserId()),
		ExpiresIn:   int32(c.Config.Jwt.ExpiresTime),
	})
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *core.Context) error {
	in := &web.AuthForgetRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.Ctx(), entity.SmsForgetAccountChannel, in.Mobile, in.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.UserService.Forget(&service.UserForgetOpt{
		Mobile:   in.Mobile,
		Password: in.Password,
		SmsCode:  in.SmsCode,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	c.SmsService.Delete(ctx.Ctx(), entity.SmsForgetAccountChannel, in.Mobile)

	return ctx.Success(&web.AuthForgetResponse{})
}

func (c *Auth) token(uid int) string {

	expiresAt := time.Now().Add(time.Second * time.Duration(c.Config.Jwt.ExpiresTime))

	// 生成登录凭证
	token := jwt.GenerateToken("api", c.Config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(uid),
		Issuer:    "im.web",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token
}

// 设置黑名单
func (c *Auth) toBlackList(ctx *core.Context) {

	session := ctx.JwtSession()
	if session != nil {
		if ex := session.ExpiresAt - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.Ctx(), session.Token, time.Duration(ex)*time.Second)
		}
	}
}
