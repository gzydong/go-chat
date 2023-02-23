package ichat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlerFunc(fn func(ctx *Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(New(c)); err != nil {
			c.JSON(http.StatusInternalServerError, &Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Meta:    initMeta(),
			})
			c.Abort()
			return
		}
	}
}
