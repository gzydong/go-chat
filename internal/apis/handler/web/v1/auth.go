package v1

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"

	"go-chat/api/pb/queue/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
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
	Rsa                 rsautil.IRsa
}

// Login 登录接口
func (c *Auth) Login(ctx *core.Context) error {
	in := &web.AuthLoginRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	user, err := c.UserService.Login(ctx.GetContext(), in.Mobile, string(password))
	if err != nil {
		return ctx.Error(err)
	}

	data := jsonutil.Marshal(queue.UserLoginRequest{
		UserId:   int32(user.Id),
		IpAddr:   ctx.Context.ClientIP(),
		Platform: in.Platform,
		Agent:    ctx.Context.GetHeader("user-agent"),
		LoginAt:  time.Now().Format(time.DateTime),
	})

	if err := c.Redis.Publish(ctx.GetContext(), entity.LoginTopic, data).Err(); err != nil {
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
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !utils.IsMobile(in.Mobile) {
		return ctx.InvalidParams("Mobile 格式错误")
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.GetContext(), entity.SmsRegisterChannel, in.Mobile, in.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if _, err := c.UserService.Register(ctx.GetContext(), &service.UserRegisterOpt{
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Password: string(password),
		Platform: in.Platform,
	}); err != nil {
		return ctx.Error(err)
	}

	c.SmsService.Delete(ctx.GetContext(), entity.SmsRegisterChannel, in.Mobile)

	return ctx.Success(&web.AuthRegisterResponse{})
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *core.Context) error {
	token := middleware.GetAuthToken(ctx.Context)

	claims, err := jwtutil.ParseWithClaims[entity.WebClaims]([]byte(c.Config.Jwt.Secret), token)
	if err == nil {
		if ex := claims.ExpiresAt.Unix() - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.GetContext(), token, time.Duration(ex)*time.Second)
		}
	}

	return ctx.Success(nil)
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *core.Context) error {
	in := &web.AuthForgetRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !utils.IsMobile(in.Mobile) {
		return ctx.InvalidParams("Mobile 格式错误")
	}

	// 验证短信验证码是否正确
	if !c.SmsService.Verify(ctx.GetContext(), entity.SmsForgetAccountChannel, in.Mobile, in.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	password, err := c.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if _, err := c.UserService.Forget(ctx.GetContext(), &service.UserForgetOpt{
		Mobile:   in.Mobile,
		Password: string(password),
		SmsCode:  in.SmsCode,
	}); err != nil {
		return ctx.Error(err)
	}

	c.SmsService.Delete(ctx.GetContext(), entity.SmsForgetAccountChannel, in.Mobile)

	return ctx.Success(&web.AuthForgetResponse{})
}

func (c *Auth) token(uid int) string {
	token, err := jwtutil.NewTokenWithClaims(
		[]byte(c.Config.Jwt.Secret), entity.WebClaims{
			UserId: int32(uid),
		},
		func(c *jwt.RegisteredClaims) {
			c.Issuer = entity.JwtIssuerWeb
		},
		jwtutil.WithTokenExpiresAt(time.Duration(c.Config.Jwt.ExpiresTime)*time.Second),
	)

	if err != nil {
		return ""
	}

	return token
}
