package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"net/http"
)

type AuthController struct {
}

// 绑定 JSON
type Login struct {
	User     string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login 登录接口
func (a *AuthController) Login(c *gin.Context) {
	//var json Login
	//if err := c.ShouldBind(&json); err != nil {
	//	c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	//	return
	//}

	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Printf("username: %s; password: %s;\n", username, password)

	// 生成登录凭证
	token, err := helper.GenerateJwtToken("api", 2054)
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
func (a *AuthController) Register(c *gin.Context) {

}

// Logout 注销接口
func (a *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// Refresh Token 刷新接口
func (a *AuthController) Refresh(c *gin.Context) {

}

// Forget 账号找回
func (a *AuthController) Forget(c *gin.Context) {

}

// SmsCode 发送短信验证码
func (a *AuthController) SmsCode(c *gin.Context) {

}
