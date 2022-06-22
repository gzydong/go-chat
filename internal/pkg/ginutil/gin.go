package ginutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func ShouldBindJSON(ctx *gin.Context, in interface{}) error {

	value, ok := in.(proto.Message)
	if !ok {
		return fmt.Errorf("no proto.Message")
	}

	if err := ctx.ShouldBindJSON(value); err != nil {
		return err
	}

	if v, ok := value.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func HandlerFunc(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = handler(c)
	}
}
