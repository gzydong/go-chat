package v1

import (
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Common struct {
	Config      *config.Config
	UsersRepo   *repo.Users
	SmsService  service.ISmsService
	UserService service.IUserService
}

// SmsCode 发送短信验证码
func (c *Common) SmsCode(ctx *core.Context) error {
	in := &web.CommonSendSmsRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	switch in.Channel {
	// 需要判断账号是否存在
	case entity.SmsLoginChannel, entity.SmsForgetAccountChannel:
		if !c.UsersRepo.IsMobileExist(ctx.GetContext(), in.Mobile) {
			return ctx.Error(entity.ErrAccountOrPassword)
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if c.UsersRepo.IsMobileExist(ctx.GetContext(), in.Mobile) {
			return ctx.Error(entity.ErrPhoneExist)
		}

	default:
		return ctx.InvalidParams("渠道不存在")
	}

	// 发送短信验证码
	code, err := c.SmsService.Send(ctx.GetContext(), in.Channel, in.Mobile)
	if err != nil {
		return ctx.Error(err)
	}

	if in.Channel == entity.SmsRegisterChannel || in.Channel == entity.SmsChangeAccountChannel {
		return ctx.Success(map[string]any{
			"is_debug": true,
			"sms_code": code,
		})
	}

	return ctx.Success(&web.CommonSendSmsResponse{})
}

// EmailCode 发送邮件验证码
func (c *Common) EmailCode(ctx *core.Context) error {

	params := &web.CommonSendEmailRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	return ctx.Success(nil)
}

// Setting 公共设置
func (c *Common) Setting(ctx *core.Context) error {
	return nil
}
