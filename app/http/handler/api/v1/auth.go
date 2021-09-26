package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/service"
	"go-chat/config"
)

type Auth struct {
	Conf        *config.Config
	UserService *service.UserService
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

	_, err := a.UserService.Register(param)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, gin.H{}, "账号注册成功")
}

// Logout 注销接口
func (a *Auth) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// Refresh Token 刷新接口
func (a *Auth) Refresh(c *gin.Context) {

}

// Forget 账号找回
func (a *Auth) Forget(c *gin.Context) {

}

// SmsCode 发送短信验证码
func (a *Auth) SmsCode(c *gin.Context) {

}
