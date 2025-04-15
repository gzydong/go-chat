package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/core/validator"
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

func (c *Context) ShouldBindProto(in any) error {
	value, ok := in.(proto.Message)
	if err := c.Context.ShouldBind(in); err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if v, ok := value.(interface{ Validate() error }); ok {
		return v.Validate()
	}

	return nil
}

// Unauthorized 未认证
func (c *Context) Unauthorized(message string) error {
	c.Context.AbortWithStatusJSON(http.StatusUnauthorized, &Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	})

	return nil
}

// Forbidden 未授权
func (c *Context) Forbidden(message string) error {
	c.Context.AbortWithStatusJSON(http.StatusForbidden, &Response{
		Code:    http.StatusForbidden,
		Message: message,
	})

	return nil
}

// InvalidParams 参数错误
func (c *Context) InvalidParams(message any) error {
	resp := &Response{Code: 400, Message: "invalid params"}

	switch msg := message.(type) {
	case error:
		resp.Message = validator.Translate(msg)
	case string:
		resp.Message = msg
	default:
		resp.Message = fmt.Sprintf("%v", msg)
	}

	c.Context.AbortWithStatusJSON(http.StatusBadRequest, resp)

	return nil
}

// Error 错误信息响应
func (c *Context) Error(err error) error {
	resp := &Response{Code: 400, Message: err.Error()}

	var e *errorx.Error
	if errors.As(err, &e) {
		resp.Code = e.Code
		resp.Message = e.Message

		if slices.Contains([]int{404, 403, 429, 400}, resp.Code) {
			c.Context.AbortWithStatusJSON(resp.Code, resp)
		} else {
			c.Context.AbortWithStatusJSON(http.StatusBadRequest, resp)
		}
	} else {
		resp.Code = 500
		resp.Message = err.Error()
		c.Context.AbortWithStatusJSON(http.StatusInternalServerError, resp)
	}

	return nil
}

// Success 成功响应(Json 数据)
func (c *Context) Success(data any, message ...string) error {
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

		var body map[string]any
		_ = json.Unmarshal(bt, &body)
		resp.Data = body
	}

	c.Context.AbortWithStatusJSON(http.StatusOK, resp)
	return nil
}

// Raw 成功响应(原始数据)
func (c *Context) Raw(value string) error {
	c.Context.Abort()
	c.Context.String(http.StatusOK, value)
	return nil
}

// AuthId 返回JWT登录用户的ID
func (c *Context) AuthId() int {
	id, ok := c.Context.Get(middleware.JWTAuthID)
	if ok {
		return id.(int)
	}

	return 0
}

// IsGuest 是否是游客(未登录状态)
func (c *Context) IsGuest() bool {
	return c.AuthId() == 0
}

func (c *Context) GetContext() context.Context {
	return c.Context.Request.Context()
}
