package response

import (
	"fmt"
	"net/http"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/validation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
)

// MarshalOptions is a configurable JSON format marshaller.
var MarshalOptions = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

// Response 返回数据结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

func NewError(c *gin.Context, code int, message ...interface{}) {
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
		if errorMessage, ok := entity.CodeMessageMap[code]; ok {
			msg = errorMessage
		} else {
			msg = entity.CodeMessageMap[entity.CodeSystemError]
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

// Success 响应成功数据
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	c.Abort()

	if value, ok := data.(proto.Message); ok {
		var val interface{}

		bt, _ := MarshalOptions.Marshal(value)
		if err := jsonutil.Decode(string(bt), &val); err != nil {
			SystemError(c, err.Error())
			return
		}

		c.JSON(http.StatusOK, &Response{Code: entity.CodeSuccess, Data: val, Message: msg})
	} else {
		c.JSON(http.StatusOK, &Response{Code: entity.CodeSuccess, Data: data, Message: msg})
	}
}

// SuccessPaginate 响应分页数据
func SuccessPaginate(c *gin.Context, rows interface{}, page, size, total int) {
	c.JSON(http.StatusOK, &Response{Code: entity.CodeSuccess, Message: "success", Data: PaginateResponse{
		Rows:     rows,
		Paginate: Paginate{Page: page, Size: size, Total: total},
	}})
	c.Abort()
}
