package v1

import (
	"context"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
)

var _ web.ICommonHandler = (*Common)(nil)

type Common struct {
	Config      *config.Config
	UsersRepo   *repo.Users
	SmsService  service.ISmsService
	UserService service.IUserService
}

func (c *Common) SendSms(ctx context.Context, in *web.CommonSendSmsRequest) (*web.CommonSendSmsResponse, error) {
	switch in.Channel {
	// 需要判断账号是否存在
	case entity.SmsLoginChannel, entity.SmsForgetAccountChannel:
		if !c.UsersRepo.IsMobileExist(ctx, in.Mobile) {
			return nil, entity.ErrAccountOrPassword
		}

	// 需要判断账号是否存在
	case entity.SmsRegisterChannel, entity.SmsChangeAccountChannel:
		if c.UsersRepo.IsMobileExist(ctx, in.Mobile) {
			return nil, entity.ErrPhoneExist
		}
	case entity.SmsOauthBindChannel:
	default:
		return nil, entity.ErrSmsChannelInvalid
	}

	// 发送短信验证码
	code, err := c.SmsService.Send(ctx, in.Channel, in.Mobile)
	if err != nil {
		return nil, err
	}

	if in.Channel == entity.SmsRegisterChannel || in.Channel == entity.SmsChangeAccountChannel || in.Channel == entity.SmsOauthBindChannel {
		return &web.CommonSendSmsResponse{
			SmsCode: code,
		}, nil
	}

	return &web.CommonSendSmsResponse{}, nil
}

func (c *Common) SendEmail(ctx context.Context, req *web.CommonSendEmailRequest) (*web.CommonSendEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Common) Test(ctx context.Context, req *web.CommonSendTestRequest) (*web.CommonSendTestResponse, error) {
	//TODO implement me
	panic("implement me")
}
