package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/encrypt"
	"go-chat/app/service"
)

type User struct {
	service    *service.UserService
	smsService *service.SmsService
}

func NewUserHandler(
	userService *service.UserService,
	smsService *service.SmsService,
) *User {
	return &User{
		service:    userService,
		smsService: smsService,
	}
}

// Detail 个人用户信息
func (u *User) Detail(ctx *gin.Context) {
	user, _ := u.service.Dao().FindById(auth.GetAuthUserID(ctx))

	response.Success(ctx, user)
}

// Setting 用户设置
func (u *User) Setting(ctx *gin.Context) {
	user, _ := u.service.Dao().FindById(auth.GetAuthUserID(ctx))

	response.Success(ctx, gin.H{
		"user_info": gin.H{
			"uid":      user.Id,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"motto":    user.Motto,
			"gender":   user.Gender,
		},
		"setting": gin.H{
			"theme_mode":            "",
			"theme_bag_img":         "",
			"theme_color":           "",
			"notify_cue_tone":       "",
			"keyboard_event_notify": "",
		},
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *gin.Context) {
	params := &request.ChangeDetailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	_, _ = u.service.Dao().BaseUpdate(&model.Users{}, entity.Map{
		"id": auth.GetAuthUserID(ctx),
	}, entity.Map{
		"nickname": params.Nickname,
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Motto,
	})

	response.Success(ctx, nil, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *gin.Context) {
	params := &request.ChangePasswordRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := u.service.UpdatePassword(auth.GetAuthUserID(ctx), params.OldPassword, params.NewPassword); err != nil {
		response.BusinessError(ctx, "密码修改失败！")
		return
	}

	response.Success(ctx, nil, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *gin.Context) {
	params := &request.ChangeMobileRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !u.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		response.BusinessError(ctx, "短信验证码填写错误！")
		return
	}

	user, _ := u.service.Dao().FindById(auth.GetAuthUserID(ctx))

	if user.Mobile != params.Mobile {
		response.BusinessError(ctx, "手机号与原手机号一致无需修改！")
		return
	}

	if !encrypt.VerifyPassword(user.Password, params.Password) {
		response.BusinessError(ctx, "账号密码填写错误！")
		return
	}

	_, err := u.service.Dao().BaseUpdate(&model.Users{}, entity.Map{"id": user.Id}, entity.Map{"mobile": params.Mobile})
	if err != nil {
		response.BusinessError(ctx, "手机号修改失败！")
		return
	}

	response.Success(ctx, nil, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *gin.Context) {
	params := &request.ChangeEmailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// todo 1.验证邮件激活码是否正确
}
