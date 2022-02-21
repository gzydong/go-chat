package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-chat/config"
)

// Cors 处理跨域请求,支持options访问
func Cors(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", conf.Cors.Origin)
		c.Header("Access-Control-Allow-Headers", conf.Cors.Headers)
		c.Header("Access-Control-Allow-Methods", conf.Cors.Methods)
		c.Header("Access-Control-Allow-Credentials", conf.Cors.Credentials)
		c.Header("Access-Control-Max-Age", conf.Cors.MaxAge)

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}
