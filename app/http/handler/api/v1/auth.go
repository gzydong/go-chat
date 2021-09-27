package v1

import (
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
	Conf        *config.Config
	UserService *service.UserService
	SmsService  *service.SmsService
	AuthToken   *cache.AuthToken
}

// Login 登录接口
func (a *Auth) Login(c *gin.Context) {
	param := &request.LoginRequest{}
	if err := c.Bind(param); err != nil {
		response.InvalidParams(c, err)
		return
	}

	user, err := a.UserService.Login(param.Username, param.Password)
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
	param := &request.RegisterRequest{}
	if err := c.Bind(param); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 验证短信验证码是否正确
	isTrue := a.SmsService.CheckSmsCode(c.Request.Context(), "register", param.Mobile, param.SmsCode)
	if !isTrue {
		response.InvalidParams(c, "短信验证码填写错误！")
	}

	_, err := a.UserService.Register(param)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, gin.H{}, "账号注册成功")
}

// Logout 退出登录接口
func (a *Auth) Logout(c *gin.Context) {
	token := helper.GetAuthToken(c)

	claims, err := helper.ParseJwtToken(a.Conf, token)
	if err != nil {
		response.Success(c, "")
		return
	}

	// 计算过期时间
	expiresAt := claims.ExpiresAt - time.Now().Unix()

	// 将 token 加入黑名单
	_ = a.AuthToken.SetBlackList(c.Request.Context(), token, int(expiresAt))

	response.Success(c, "", "退出成功！")
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(c *gin.Context) {
	token, e := helper.GenerateJwtToken(a.Conf, "api", c.GetInt("__user_id__"))
	if e != nil {
		response.BusinessError(c, "Token 刷新失败，请稍后再试!")
		return
	}

	// todo 将之前的 token 加入黑名单

	response.Success(c, map[string]interface{}{
		"type":       "Bearer",
		"token":      token["token"],
		"expires_in": token["expired_at"],
	})
}

// Forget 账号找回接口
func (a *Auth) Forget(c *gin.Context) {
	param := &request.ForgetRequest{}

	if err := c.Bind(param); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 验证短信验证码是否正确
	isTrue := a.SmsService.CheckSmsCode(c.Request.Context(), "forget", param.Mobile, param.SmsCode)
	if !isTrue {
		response.InvalidParams(c, "短信验证码填写错误！")
	}

	_, err := a.UserService.Forget(param)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, gin.H{}, "账号成功找回")
}

// SmsCode 发送短信验证码
func (a *Auth) SmsCode(c *gin.Context) {
	param := &request.SmsCodeRequest{}

	if err := c.Bind(param); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// 发送短信验证码
	if err := a.SmsService.SendSmsCode(c.Request.Context(), param.Channel, param.Mobile); err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, gin.H{}, "发送成功！")
}
