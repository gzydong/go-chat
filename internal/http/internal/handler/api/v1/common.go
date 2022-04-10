package v1

import (
	"github.com/gin-gonic/gin"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/email"
	"go-chat/internal/service"
)

type Common struct {
	config      *config.Config
	smsService  *service.SmsService
	userService *service.UserService
}

func NewCommonHandler(
	config *config.Config,
	sms *service.SmsService,
	userService *service.UserService,
) *Common {
	return &Common{
		config:      config,
		smsService:  sms,
		userService: userService,
	}
}

// SmsCode 发送短信验证码
func (c *Common) SmsCode(ctx *gin.Context) {
	params := &request.SmsCodeRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	switch params.Channel {
	// 需要判断账号是否存在
	case entity.SmsLoginChannel, entity.SmsForgetAccountChannel:
		if !c.userService.Dao().IsMobileExist(params.Mobile) {
			response.BusinessError(ctx, "账号不存在或密码错误！")
			return
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if c.userService.Dao().IsMobileExist(params.Mobile) {
			response.BusinessError(ctx, "手机号已被他人使用！")
			return
		}
	default:
		response.BusinessError(ctx, "发送异常！")
		return
	}

	// 发送短信验证码
	code, err := c.smsService.SendSmsCode(ctx.Request.Context(), params.Channel, params.Mobile)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if params.Channel == entity.SmsRegisterChannel {
		response.Success(ctx, entity.MapStrAny{
			"is_debug": true,
			"sms_code": code,
		}, "发送成功！")
	} else {
		response.Success(ctx, nil, "发送成功！")
	}
}

// EmailCode 发送邮件验证码
func (c *Common) EmailCode(ctx *gin.Context) {
	params := &request.EmailCodeRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	data := &email.Options{
		To:      []string{params.Email},
		Subject: "Lumen IM(绑定邮箱验证码)",
		Body:    "xxxxxxx",
	}

	// todo 需修改
	go func() {
		_ = email.SendMail(c.config.Email, data)
	}()

	response.Success(ctx, nil, "发送成功！")
}

// Setting 公共设置
func (c *Common) Setting(ctx *gin.Context) {

}
