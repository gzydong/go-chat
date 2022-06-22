package ginutil

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/validation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	CodeSuccess            = 200 // 状态成功
	CodeInvalidParamsError = 305 // 参数错误
	CodeBusinessError      = 400 // 业务错误
	CodeNotLoginError      = 401 // 未登录
	CodeUnauthorizedError  = 403 // 未授权
	CodeSystemError        = 500 // 服务器异常
)

var (
	// CodeMessageMap 错误码对应消息
	CodeMessageMap = map[int]string{
		CodeInvalidParamsError: "参数错误",
		CodeUnauthorizedError:  "未授权",
		CodeNotLoginError:      "未登录",
		CodeBusinessError:      "业务错误",
		CodeSystemError:        "系统错误，请重试",
	}
)

// MarshalOptions is a configurable JSON format marshaller.
var MarshalOptions = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(c *gin.Context, code int, message ...interface{}) error {
	// 判断错误类型
	var msg string
	if message != nil {
		switch message[0].(type) {
		case error:
			msg = validation.Translate(message[0].(error))
		case string:
			msg = message[0].(string)
		default:
			msg = fmt.Sprintf("%v", message[0])
		}
	} else {
		if errorMessage, ok := CodeMessageMap[code]; ok {
			msg = errorMessage
		} else {
			msg = CodeMessageMap[CodeSystemError]
		}
	}

	status := http.StatusOK
	if code == 404 {
		status = http.StatusNotFound
	} else if code == 401 {
		status = http.StatusUnauthorized
	} else if code == 403 {
		status = http.StatusForbidden
	}

	c.JSON(status, &Response{Code: code, Message: msg})

	return nil
}

func InvalidParams(c *gin.Context, message ...interface{}) error {
	return NewError(c, CodeInvalidParamsError, message...)
}

func Unauthorized(c *gin.Context, message ...interface{}) error {
	return NewError(c, CodeUnauthorizedError, message...)
}

func NotLogin(c *gin.Context, message ...interface{}) error {
	return NewError(c, CodeNotLoginError, message...)
}

func BusinessError(c *gin.Context, message ...interface{}) error {
	return NewError(c, CodeBusinessError, message...)
}

func SystemError(c *gin.Context, message ...interface{}) error {
	return NewError(c, CodeSystemError, message...)
}

// Success 响应成功数据
func Success(c *gin.Context, data interface{}, message ...string) error {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	c.Abort()

	if value, ok := data.(proto.Message); ok {
		var val interface{}

		bt, _ := MarshalOptions.Marshal(value)
		if err := jsonutil.Decode(string(bt), &val); err != nil {
			return SystemError(c, err.Error())
		}

		c.JSON(http.StatusOK, &Response{Code: CodeSuccess, Data: val, Message: msg})
	} else {
		c.JSON(http.StatusOK, &Response{Code: CodeSuccess, Data: data, Message: msg})
	}

	return nil
}

type Paginate struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type PaginateResponse struct {
	Rows     interface{} `json:"rows"`
	Paginate Paginate    `json:"paginate"`
}

// SuccessPaginate 响应分页数据
func SuccessPaginate(c *gin.Context, rows interface{}, page, size, total int) error {
	c.JSON(http.StatusOK, &Response{Code: entity.CodeSuccess, Message: "success", Data: PaginateResponse{
		Rows:     rows,
		Paginate: Paginate{Page: page, Size: size, Total: total},
	}})
	c.Abort()

	return nil
}
