package ichat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/validation"
	"google.golang.org/protobuf/proto"
)

type Context struct {
	Context *gin.Context
	data    map[string]interface{}
}

// WithData 挂载错误详情信息
func (c *Context) WithData(data map[string]interface{}) *Context {
	c.data = data
	return c
}

// InvalidParams 参数错误
func (c *Context) InvalidParams(value interface{}) error {
	c.Context.Abort()

	var msg string

	switch value.(type) {
	case error:
		msg = validation.Translate(value.(error))
	case string:
		msg = value.(string)
	default:
		msg = fmt.Sprintf("%v", value)
	}

	c.Context.JSON(http.StatusOK, gin.H{
		"code":    305,
		"message": msg,
	})

	return nil
}

// BusinessError 业务错误
func (c *Context) BusinessError(error string) error {
	c.Context.Abort()

	c.Context.JSON(http.StatusOK, gin.H{
		"code":    400,
		"message": error,
	})

	return nil
}

// Error 系统错误
func (c *Context) Error(error string) error {
	c.Context.Abort()
	c.Context.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": error,
	})

	return nil
}

// Success 成功响应(Json 数据)
func (c *Context) Success(data interface{}, message ...string) error {

	resp := make(map[string]interface{})
	resp["code"] = 200
	resp["message"] = "success"

	if len(message) > 0 {
		resp["message"] = message[0]
	}

	if value, ok := data.(proto.Message); ok {
		var body interface{}

		bt, _ := MarshalOptions.Marshal(value)
		if err := jsonutil.Decode(string(bt), &body); err != nil {
			return c.Error(err.Error())
		}

		resp["data"] = body
	} else {
		if data != nil {
			resp["data"] = data
		}
	}

	c.Context.Abort()
	c.Context.JSON(http.StatusOK, &resp)

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
	return 0
}
