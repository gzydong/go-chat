package entity

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
