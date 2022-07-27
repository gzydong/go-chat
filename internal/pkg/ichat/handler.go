package ichat

import (
	"github.com/gin-gonic/gin"
)

func HandlerFunc(fn func(ctx *Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = fn(New(c))
	}
}
