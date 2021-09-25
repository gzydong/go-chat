package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/config"
)

type Auth struct {
	Conf *config.Config
}

// 绑定 JSON
type Login struct {
	User     string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login 登录接口
func (a *Auth) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Printf("username: %s; password: %s;\n", username, password)

	// 生成登录凭证
	token, err := helper.GenerateJwtToken(a.Conf, "api", 2054)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    305,
			"message": "登录失败，请稍后再试！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"type":       "Bearer",
			"token":      token["token"],
			"expires_in": token["expired_at"],
		},
	})
}

// Register 注册接口
func (a *Auth) Register(c *gin.Context) {

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
