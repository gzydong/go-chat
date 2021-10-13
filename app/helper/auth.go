package helper

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
)

// GetAuthUserID 获取授权登录的用户ID
func GetAuthUserID(c *gin.Context) int {
	return c.GetInt(entity.LoginUserID)
}
