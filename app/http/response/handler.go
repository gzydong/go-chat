package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) *Response

func Handler(fn HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, fn(context))
	}
}
