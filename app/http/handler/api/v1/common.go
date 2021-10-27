package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/email"
	"go-chat/app/service"
	"go-chat/config"
)

type Common struct {
	config     *config.Config
	smsService *service.SmsService
	userRepo   *dao.UserDao
}

func NewCommonHandler(
	config *config.Config,
	sms *service.SmsService,
	user *dao.UserDao,
) *Common {
	return &Common{
		config:     config,
		smsService: sms,
		userRepo:   user,
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
		if !c.userRepo.IsMobileExist(params.Mobile) {
			response.BusinessError(ctx, "账号不存在！")
			return
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if c.userRepo.IsMobileExist(params.Mobile) {
			response.BusinessError(ctx, "手机号已被他人使用！")
			return
		}
	}

	// 发送短信验证码
	if err := c.smsService.SendSmsCode(ctx.Request.Context(), params.Channel, params.Mobile); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil, "发送成功！")
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
		Body:    "",
	}

	go func() {
		_ = email.SendMail(c.config.Email, data)
	}()

	response.Success(ctx, nil, "发送成功！")
}

// Setting 公共设置
func (c *Common) Setting(ctx *gin.Context) {

}
