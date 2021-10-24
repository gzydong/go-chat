package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/validator"
)

// Response 返回数据结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(c *gin.Context, code int, message ...interface{}) {
	// 判断错误类型
	var msg string
	if message != nil {
		switch message[0].(type) {
		case error:
			msg = validator.Translate(message[0].(error))
		case string:
			msg = message[0].(string)
		default:
			msg = fmt.Sprintf("%v", message[0])
		}
	} else {
		if errorMessage, ok := entity.CodeMessageMap[code]; ok {
			msg = errorMessage
		} else {
			msg = entity.CodeMessageMap[entity.CodeSystemError]
		}
	}

	status := http.StatusOK
	if code == 404 {
		status = http.StatusNotFound
	}

	c.JSON(status, &Response{Code: code, Message: msg})
}

func InvalidParams(c *gin.Context, message ...interface{}) {
	NewError(c, entity.CodeInvalidParamsError, message...)
}

func Unauthorized(c *gin.Context, message ...interface{}) {
	NewError(c, entity.CodeUnauthorizedError, message...)
}

func NotLogin(c *gin.Context, message ...interface{}) {
	NewError(c, entity.CodeNotLoginError, message...)
}

func BusinessError(c *gin.Context, message ...interface{}) {
	NewError(c, entity.CodeBusinessError, message...)
}

func SystemError(c *gin.Context, message ...interface{}) {
	NewError(c, entity.CodeSystemError, message...)
}

func Success(c *gin.Context, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, &Response{Code: entity.CodeSuccess, Data: data, Message: msg})
	c.Abort()
}
