package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", config.GlobalConfig.Cors.Origin)
		c.Header("Access-Control-Allow-Headers", config.GlobalConfig.Cors.Headers)
		c.Header("Access-Control-Allow-Methods", config.GlobalConfig.Cors.Methods)
		c.Header("Access-Control-Allow-Credentials", config.GlobalConfig.Cors.Credentials)
		c.Header("Access-Control-Max-Age", config.GlobalConfig.Cors.MaxAge)

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}
