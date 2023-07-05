package admin

import "github.com/gin-gonic/gin"

func Permission() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
