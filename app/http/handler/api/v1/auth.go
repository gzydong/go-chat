package v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/helper"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/service"
	"go-chat/config"
)

type Auth struct {
	config         *config.Config
	userService    *service.UserService
	smsService     *service.SmsService
	authTokenCache *cache.AuthTokenCache
	redisLock      *cache.RedisLock
}

func NewAuthHandler(
	config *config.Config,
	userService *service.UserService,
	smsService *service.SmsService,
	tokenCache *cache.AuthTokenCache,
	lock *cache.RedisLock,
) *Auth {
	return &Auth{
		config:         config,
		userService:    userService,
		smsService:     smsService,
		authTokenCache: tokenCache,
		redisLock:      lock,
	}
}

// Login 登录接口
func (a *Auth) Login(ctx *gin.Context) {
	params := &request.LoginRequest{}
	if err := ctx.Bind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	user, err := a.userService.Login(params.Mobile, params.Password)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 生成登录凭证
	token, e := helper.GenerateJwtToken(a.config, "api", user.ID)
	if e != nil {
		response.BusinessError(ctx, "登录失败，请稍后再试")
		return
	}

	response.Success(ctx, map[string]interface{}{
		"type":       "Bearer",
		"token":      token["token"],
		"expires_in": token["expired_at"],
	})
}

// Register 注册接口
func (a *Auth) Register(ctx *gin.Context) {
	params := &request.RegisterRequest{}
	if err := ctx.Bind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	_, err := a.userService.Register(params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	a.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile)

	response.Success(ctx, gin.H{}, "账号注册成功")
}

// Logout 退出登录接口
func (a *Auth) Logout(ctx *gin.Context) {
	token := helper.GetAuthToken(ctx)

	claims, err := helper.ParseJwtToken(a.config.Jwt.Secret, token)
	if err != nil {
		response.Success(ctx, gin.H{})
		return
	}

	// 计算过期时间
	expiresAt := claims.ExpiresAt - time.Now().Unix()

	// 将 token 加入黑名单
	_ = a.authTokenCache.SetBlackList(ctx.Request.Context(), token, int(expiresAt))

	response.Success(ctx, gin.H{}, "退出成功！")
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(ctx *gin.Context) {
	token, err := helper.GenerateJwtToken(a.config, "api", helper.GetAuthUserID(ctx))
	if err != nil {
		response.BusinessError(ctx, "Token 刷新失败，请稍后再试!")
		return
	}

	// todo 将之前的 token 加入黑名单

	response.Success(ctx, gin.H{
		"type":       "Bearer",
		"token":      token["token"],
		"expires_in": token["expired_at"],
	})
}

// Forget 账号找回接口
func (a *Auth) Forget(ctx *gin.Context) {
	params := &request.ForgetRequest{}

	if err := ctx.Bind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	_, err := a.userService.Forget(params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	a.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile)

	response.Success(ctx, gin.H{}, "账号成功找回")
}
