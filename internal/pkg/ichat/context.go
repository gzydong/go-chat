package ichat

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/ichat/middleware"
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
}

func New(ctx *gin.Context) *Context {
	return &Context{ctx}
}

// Unauthorized 未认证
func (c *Context) Unauthorized(message string) error {

	c.Context.Abort()
	c.Context.JSON(http.StatusUnauthorized, &Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	})

	return nil
}

// Forbidden 未授权
func (c *Context) Forbidden(message string) error {

	c.Context.Abort()
	c.Context.JSON(http.StatusForbidden, &Response{
		Code:    http.StatusForbidden,
		Message: message,
	})

	return nil
}

// InvalidParams 参数错误
func (c *Context) InvalidParams(message interface{}) error {

	resp := &Response{Code: 305, Message: "invalid params"}

	switch msg := message.(type) {
	case error:
		resp.Message = validation.Translate(msg)
	case string:
		resp.Message = msg
	default:
		resp.Message = fmt.Sprintf("%v", msg)
	}

	c.Context.Abort()
	c.Context.JSON(http.StatusOK, resp)

	return nil
}

// ErrorBusiness 业务错误
func (c *Context) ErrorBusiness(message interface{}) error {

	resp := &Response{Code: 400, Message: "business error"}

	switch msg := message.(type) {
	case error:
		resp.Message = msg.Error()
	case string:
		resp.Message = msg
	default:
		resp.Message = fmt.Sprintf("%v", msg)
	}

	meta := make(map[string]interface{})
	_, _, line, ok := runtime.Caller(1)
	if ok {
		meta["error_line"] = line
	}

	resp.Meta = meta

	c.Context.Abort()
	c.Context.JSON(http.StatusOK, resp)

	return nil
}

// Error 系统错误
func (c *Context) Error(error string) error {

	meta := make(map[string]interface{})
	_, _, line, ok := runtime.Caller(1)
	if ok {
		meta["error_line"] = line
	}

	c.Context.Abort()
	c.Context.JSON(http.StatusInternalServerError, &Response{
		Code:    500,
		Message: error,
		Meta:    meta,
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

	// 检测是否是 proto 对象
	if value, ok := data.(proto.Message); ok {
		bt, _ := MarshalOptions.Marshal(value)

		var data interface{}
		if err := jsonutil.Decode(string(bt), &data); err != nil {
			return c.Error(err.Error())
		}

		resp.Data = data
	}

	c.Context.Abort()
	c.Context.JSON(http.StatusOK, resp)

	return nil
}

// Raw 成功响应(原始数据)
func (c *Context) Raw(value string) error {

	c.Context.Abort()
	c.Context.String(http.StatusOK, value)

	return nil
}

// UserId 返回登录用户的UID
func (c *Context) UserId() int {

	if session := c.JwtSession(); session != nil {
		return session.Uid
	}

	return 0
}

// JwtSession 返回登录用户的JSession
func (c *Context) JwtSession() *middleware.JSession {

	data, isOk := c.Context.Get(middleware.JWTSessionConst)
	if !isOk {
		return nil
	}

	return data.(*middleware.JSession)
}

// IsGuest 是否是游客(未登录状态)
func (c *Context) IsGuest() bool {
	return c.UserId() == 0
}

func (c *Context) Ctx() context.Context {
	return c.Context.Request.Context()
}
