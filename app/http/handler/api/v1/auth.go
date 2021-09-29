package v1

import (
	"go-chat/app/entity"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/helper"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/service"
	"go-chat/config"
)

type Auth struct {
	Conf           *config.Config
	UserService    *service.UserService
	SmsService     *service.SmsService
	AuthTokenCache *cache.AuthTokenCache
	RedisLock      *cache.RedisLock
}

// Login 登录接口
func (a *Auth) Login(c *gin.Context) {
	params := &request.LoginRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	user, err := a.UserService.Login(params.Mobile, params.Password)
	if err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 生成登录凭证
	token, e := helper.GenerateJwtToken(a.Conf, "api", user.ID)
	if e != nil {
		response.BusinessError(c, "登录失败，请稍后再试")
		return
	}

	response.Success(c, map[string]interface{}{
		"type":       "Bearer",
		"token":      token["token"],
		"expires_in": token["expired_at"],
	})
}

// Register 注册接口
func (a *Auth) Register(c *gin.Context) {
	params := &request.RegisterRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.SmsService.CheckSmsCode(c.Request.Context(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(c, "短信验证码填写错误！")
		return
	}

	_, err := a.UserService.Register(params)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	a.SmsService.DeleteSmsCode(c.Request.Context(), entity.SmsRegisterChannel, params.Mobile)

	response.Success(c, gin.H{}, "账号注册成功")
}

// Logout 退出登录接口
func (a *Auth) Logout(c *gin.Context) {
	token := helper.GetAuthToken(c)

	claims, err := helper.ParseJwtToken(a.Conf, token)
	if err != nil {
		response.Success(c, gin.H{})
		return
	}

	// 计算过期时间
	expiresAt := claims.ExpiresAt - time.Now().Unix()

	// 将 token 加入黑名单
	_ = a.AuthTokenCache.SetBlackList(c.Request.Context(), token, int(expiresAt))

	response.Success(c, gin.H{}, "退出成功！")
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(c *gin.Context) {
	token, err := helper.GenerateJwtToken(a.Conf, "api", c.GetInt("__user_id__"))
	if err != nil {
		response.BusinessError(c, "Token 刷新失败，请稍后再试!")
		return
	}

	// todo 将之前的 token 加入黑名单

	response.Success(c, gin.H{
		"type":       "Bearer",
		"token":      token["token"],
		"expires_in": token["expired_at"],
	})
}

// Forget 账号找回接口
func (a *Auth) Forget(c *gin.Context) {
	params := &request.ForgetRequest{}

	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 验证短信验证码是否正确
	if !a.SmsService.CheckSmsCode(c.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		response.InvalidParams(c, "短信验证码填写错误！")
		return
	}

	// 密码找回
	_, err := a.UserService.Forget(params)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	a.SmsService.DeleteSmsCode(c.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile)

	response.Success(c, gin.H{}, "账号成功找回")
}
