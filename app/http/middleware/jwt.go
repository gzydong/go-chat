package middleware

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/config"
	"net/http"
)

// ApiAuth 授权中间件
// guard 授权守卫
func JwtAuth(conf *config.Config, guard string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := helper.GetAuthToken(c)
		if token == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 40001,
				"msg":  "请登录后操作!",
			})
			return
		}

		jwt, err := helper.ParseJwtToken(conf, token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 40002,
				"msg":  "Token 信息验证失败!",
			})
			return
		}

		// 判断权限认证守卫是否一致
		if jwt.Guard != guard {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 40002,
				"msg":  "Token 信息验证失败!",
			})
			return
		}

		// todo 黑名单判断

		c.Set("user_id", jwt.UserId)

		c.Next()
	}
}
