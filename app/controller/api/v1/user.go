package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
}

// GetUserDetail 获取登录用户信息
func (u *UserController) Detail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}
