package v1

import (
	"go-chat/app/pkg/auth"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/entity"
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
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	user, err := a.userService.Login(params.Mobile, params.Password)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(a.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := auth.SignJwtToken("api", a.config.Jwt.Secret, &auth.JwtOptions{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(user.ID),
	})

	response.Success(ctx, map[string]interface{}{
		"type":       "Bearer",
		"token":      token,
		"expires_in": expiresAt,
	})
}

// Register 注册接口
func (a *Auth) Register(ctx *gin.Context) {
	params := &request.RegisterRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	if _, err := a.userService.Register(params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	a.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsRegisterChannel, params.Mobile)

	response.Success(ctx, gin.H{}, "账号注册成功")
}

// Logout 退出登录接口
func (a *Auth) Logout(ctx *gin.Context) {

	info := ctx.GetStringMapString("jwt")

	expiresAt, _ := strconv.Atoi(info["expires_at"])

	// 将 token 加入黑名单
	_ = a.authTokenCache.SetBlackList(ctx.Request.Context(), info["token"], expiresAt-int(time.Now().Unix()))

	response.Success(ctx, gin.H{}, "退出成功！")
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(ctx *gin.Context) {
	expiresAt := time.Now().Add(time.Second * time.Duration(a.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := auth.SignJwtToken("api", a.config.Jwt.Secret, &auth.JwtOptions{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(auth.GetAuthUserID(ctx)),
	})

	response.Success(ctx, gin.H{
		"type":       "Bearer",
		"token":      token,
		"expires_in": expiresAt,
	})
}

// Forget 账号找回接口
func (a *Auth) Forget(ctx *gin.Context) {
	params := &request.ForgetRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(ctx, "短信验证码填写错误！")
		return
	}

	if _, err := a.userService.Forget(params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	a.smsService.DeleteSmsCode(ctx.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile)

	response.Success(ctx, gin.H{}, "账号成功找回")
}
