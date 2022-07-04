package ichat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/validation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MarshalOptions is a configurable JSON format marshaller.
var MarshalOptions = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

type Context struct {
	Context *gin.Context
	meta    interface{}
}

// WithMeta 挂载错误详情信息
func (c *Context) WithMeta(data map[string]interface{}) *Context {
	c.meta = data
	return c
}

func (c *Context) Unauthorized(message string) error {
	c.Context.Abort()

	c.Context.JSON(http.StatusForbidden, &Response{
		Code:    403,
		Message: message,
		Meta:    c.meta,
	})

	return nil
}

// InvalidParams 参数错误
func (c *Context) InvalidParams(value interface{}) error {
	c.Context.Abort()

	var msg string

	switch val := value.(type) {
	case error:
		msg = validation.Translate(val)
	case string:
		msg = val
	default:
		msg = fmt.Sprintf("%v", val)
	}

	c.Context.JSON(http.StatusOK, &Response{
		Code:    305,
		Message: msg,
		Meta:    c.meta,
	})

	return nil
}

// BusinessError 业务错误
func (c *Context) BusinessError(message interface{}) error {

	resp := &Response{
		Code:    400,
		Message: "business error",
		Meta:    c.meta,
	}

	switch msg := message.(type) {
	case error:
		resp.Message = msg.Error()
	case string:
		resp.Message = msg
	default:
		resp.Message = fmt.Sprintf("%v", msg)
	}

	c.Context.Abort()
	c.Context.JSON(http.StatusOK, resp)

	return nil
}

// Error 系统错误
func (c *Context) Error(error string) error {

	c.Context.Abort()

	c.Context.JSON(http.StatusInternalServerError, &Response{
		Code:    500,
		Message: error,
		Meta:    c.meta,
	})

	return nil
}

// Success 成功响应(Json 数据)
func (c *Context) Success(data interface{}, message ...string) error {

	resp := &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}

	if len(message) > 0 {
		resp.Message = message[0]
	}

	if value, ok := data.(proto.Message); ok {
		var body interface{}

		bt, _ := MarshalOptions.Marshal(value)
		if err := jsonutil.Decode(string(bt), &body); err != nil {
			return c.Error(err.Error())
		}

		resp.Data = body
	}

	c.Context.JSON(http.StatusOK, resp)

	return nil
}

func (c *Context) Paginate(items interface{}, page, size, total int) error {
	c.Context.Abort()

	c.Context.JSON(http.StatusOK, &Response{
		Code:    200,
		Message: "success",
		Data: &PaginateResponse{
			Items: items,
			Paginate: Paginate{
				Page:  page,
				Size:  size,
				Total: total,
			},
		},
	})

	return nil
}

// Raw 成功响应(原始数据)
func (c *Context) Raw(value string) error {
	c.Context.Abort()

	c.Context.String(http.StatusOK, value)

	return nil
}

// LoginUID 返回登录用户的UID
func (c *Context) LoginUID() int {
	return c.Context.GetInt("__UID__")
}

// IsGuest 是否是游客(未登录状态)
func (c *Context) IsGuest() bool {
	return c.LoginUID() == 0
}
