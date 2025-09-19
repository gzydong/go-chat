package router

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/core/errorx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ProtoJson = protojson.MarshalOptions{
	UseProtoNames:   true, // 使用下划线命名
	EmitUnpopulated: true, // 空值也返回
}

type Response struct{}

func (r *Response) ShouldProto(ctx *gin.Context, in any) error {
	if err := ctx.ShouldBind(in); err != nil {
		return errorx.New(400, err.Error())
	}

	if v, ok := in.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return errorx.New(400, err.Error())
		}
	}

	return nil
}

func (r *Response) ErrorResponse(c *gin.Context, err error) {
	// 这里需要断言这个错误是否是指定错误码
	var e *errorx.Error
	if errors.As(err, &e) {
		if slices.Contains([]int{404, 403, 429, 400}, e.Code) {
			c.AbortWithStatusJSON(e.Code, gin.H{"code": e.Code, "message": e.Message})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": e.Code, "message": e.Message})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
	}
}

func (r *Response) SuccessResponse(c *gin.Context, data any) {
	// 检测是否是 proto 对象
	if value, ok := data.(proto.Message); ok {
		body, _ := ProtoJson.Marshal(value)
		c.Data(http.StatusOK, "application/json", body)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
