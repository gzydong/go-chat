package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/repository"
	"go-chat/app/service"
)

type Common struct {
	SmsService *service.SmsService
	UserRepo   *repository.UserRepository
}

// SmsCode 发送短信验证码
func (a *Common) SmsCode(c *gin.Context) {
	params := &request.SmsCodeRequest{}

	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	switch params.Channel {
	// 需要判断账号是否存在
	case entity.SmsLoginChannel, entity.SmsForgetAccountChannel:
		if !a.UserRepo.IsMobileExist(params.Mobile) {
			response.BusinessError(c, "账号不存在！")
			return
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if a.UserRepo.IsMobileExist(params.Mobile) {
			response.BusinessError(c, "手机号已被他人使用！")
			return
		}
	}

	// 发送短信验证码
	if err := a.SmsService.SendSmsCode(c.Request.Context(), params.Channel, params.Mobile); err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, nil, "发送成功！")
}
