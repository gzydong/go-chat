package ichat

import (
	"github.com/gin-gonic/gin"
)

func HandlerFunc(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler(c)
	}
}

func NewHandlerFunc(handler func(ctx *Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler(&Context{Context: c})
	}
}
