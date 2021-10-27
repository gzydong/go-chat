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
	config      *config.Config
	userService *service.UserService
	smsService  *service.SmsService
	token       *cache.AuthTokenCache
	redisLock   *cache.RedisLock
}

type CertToken struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func NewAuthHandler(
	config *config.Config,
	userService *service.UserService,
	smsService *service.SmsService,
	tokenCache *cache.AuthTokenCache,
	redisLock *cache.RedisLock,
) *Auth {
	return &Auth{
		config:      config,
		userService: userService,
		smsService:  smsService,
		token:       tokenCache,
		redisLock:   redisLock,
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

	// 处理登录日志消息

	response.Success(ctx, a.createToken(user.ID))
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
	a.toBlackList(ctx)

	response.Success(ctx, gin.H{}, "退出成功！")
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(ctx *gin.Context) {
	tokenInfo := a.createToken(auth.GetAuthUserID(ctx))

	a.toBlackList(ctx)

	response.Success(ctx, tokenInfo)
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

func (a *Auth) createToken(uid int) *CertToken {
	expiresAt := time.Now().Add(time.Second * time.Duration(a.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := auth.SignJwtToken("api", a.config.Jwt.Secret, &auth.JwtOptions{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(uid),
	})

	return &CertToken{
		Type:      "Bearer",
		Token:     token,
		ExpiresIn: expiresAt,
	}
}

// 设置黑名单
func (a *Auth) toBlackList(ctx *gin.Context) {
	info := ctx.GetStringMapString("jwt")

	expiresAt, _ := strconv.Atoi(info["expires_at"])

	ex := expiresAt - int(time.Now().Unix())

	// 将 token 加入黑名单
	_ = a.token.SetBlackList(ctx.Request.Context(), info["token"], ex)
}
