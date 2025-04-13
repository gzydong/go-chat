package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ICorsOptions interface {
	GetOrigin() string
	GetHeaders() string
	GetMethods() string
	GetCredentials() string
	GetMaxAge() string
}

// Cors 处理跨域请求
func Cors(options ICorsOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := strings.ReplaceAll(uuid.NewString(), "-", "")

		c.Header("Access-Control-Allow-Origin", options.GetOrigin())
		c.Header("Access-Control-Allow-Headers", options.GetHeaders())
		c.Header("Access-Control-Allow-Methods", options.GetMethods())
		c.Header("Access-Control-Allow-Credentials", options.GetCredentials())
		c.Header("Access-Control-Max-Age", options.GetMaxAge())
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("trace-id", traceId)

		// 放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// nolint
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "trace_id", traceId))

		// 处理请求
		c.Next()
	}
}
