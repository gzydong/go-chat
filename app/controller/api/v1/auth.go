package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"net/http"
)

type AuthController struct {
}

// Login 用户登录接口
func (a *AuthController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Printf("username: %s; password: %s;", username, password)

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

// Register 用户注册接口
func (a *AuthController) Register(c *gin.Context) {

}
