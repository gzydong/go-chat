package response

import (
	"fmt"

	"go-chat/app/entity"
)

// 返回数据结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message ...interface{}) *Response {
	// 判断错误类型
	var msg string
	if message != nil {
		switch message[0].(type) {
		case error:
			msg = message[0].(error).Error()
		case string:
			msg = message[0].(string)
		default:
			msg = fmt.Sprintf("%v", message[0])
		}
	}

	// 没有消息内容
	if msg == "" {
		if msg1, ok := entity.CodeMessageMap[code]; ok {
			msg = msg1
		} else {
			msg = entity.CodeMessageMap[entity.CodeSystemError]
		}
	}

	return &Response{Code: code, Message: msg}
}

func InvalidParams(message ...interface{}) *Response {
	return NewError(entity.CodeInvalidParamsError, message...)
}

func Unauthorized(message ...interface{}) *Response {
	return NewError(entity.CodeUnauthorizedError, message...)
}

func NotLogin(message ...interface{}) *Response {
	return NewError(entity.CodeNotLoginError, message...)
}

func BusinessError(message ...interface{}) *Response {
	return NewError(entity.CodeBusinessError, message...)
}

func SystemError(message ...interface{}) *Response {
	return NewError(entity.CodeSystemError, message...)
}

func Success(data interface{}, message ...string) *Response {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	return &Response{Code: entity.CodeSuccess, Data: data, Message: msg}
}
