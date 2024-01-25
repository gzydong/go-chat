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
	UsersRepo *repo.Users

	Config      *config.Config
	SmsService  service.ISmsService
	UserService service.IUserService
}

// SmsCode 发送短信验证码
func (c *Common) SmsCode(ctx *core.Context) error {
	in := &web.CommonSendSmsRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	switch in.Channel {
	// 需要判断账号是否存在
	case entity.SmsLoginChannel, entity.SmsForgetAccountChannel:
		if !c.UsersRepo.IsMobileExist(ctx.Ctx(), in.Mobile) {
			return ctx.ErrorBusiness("账号不存在或密码错误！")
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if c.UsersRepo.IsMobileExist(ctx.Ctx(), in.Mobile) {
			return ctx.ErrorBusiness("手机号已被他人使用！")
		}

	default:
		return ctx.ErrorBusiness("发送异常！")
	}

	// 发送短信验证码
	code, err := c.SmsService.Send(ctx.Ctx(), in.Channel, in.Mobile)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
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
