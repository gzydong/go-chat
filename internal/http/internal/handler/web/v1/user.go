package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/model"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service"
	"go-chat/internal/service/organize"
)

type User struct {
	service      *service.UserService
	smsService   *service.SmsService
	organizeServ *organize.OrganizeService
}

func NewUserHandler(
	userService *service.UserService,
	smsService *service.SmsService,
	organizeServ *organize.OrganizeService,
) *User {
	return &User{
		service:      userService,
		smsService:   smsService,
		organizeServ: organizeServ,
	}
}

// Detail 个人用户信息
func (u *User) Detail(ctx *gin.Context) error {

	uid := jwtutil.GetUid(ctx)

	user, err := u.service.Dao().FindById(uid)
	if err != nil {
		return ginutil.SystemError(ctx, err.Error())
	}

	return ginutil.Success(ctx, &web.GetUserInfoResponse{
		Id:       user.Id,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Motto:    user.Motto,
		Email:    user.Email,
	})
}

// Setting 用户设置
func (u *User) Setting(ctx *gin.Context) error {
	uid := jwtutil.GetUid(ctx)

	user, _ := u.service.Dao().FindById(uid)

	isOk, _ := u.organizeServ.Dao().IsQiyeMember(uid)

	return ginutil.Success(ctx, entity.H{
		"user_info": entity.H{
			"uid":      user.Id,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"motto":    user.Motto,
			"gender":   user.Gender,
			"is_qiye":  isOk,
			"mobile":   user.Mobile,
			"email":    user.Email,
		},
		"setting": entity.H{
			"theme_mode":            "",
			"theme_bag_img":         "",
			"theme_color":           "",
			"notify_cue_tone":       "",
			"keyboard_event_notify": "",
		},
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *gin.Context) error {
	params := &web.ChangeUserDetailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	_, _ = u.service.Dao().BaseUpdate(&model.Users{}, entity.MapStrAny{
		"id": jwtutil.GetUid(ctx),
	}, entity.MapStrAny{
		"nickname": params.Nickname,
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Motto,
	})

	return ginutil.Success(ctx, nil, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *gin.Context) error {
	params := &web.ChangeUserPasswordRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)

	if uid == 2054 || uid == 2055 {
		return ginutil.BusinessError(ctx, "预览账号不支持修改密码！")
	}

	if err := u.service.UpdatePassword(jwtutil.GetUid(ctx), params.OldPassword, params.NewPassword); err != nil {
		return ginutil.BusinessError(ctx, "密码修改失败！")
	}

	return ginutil.Success(ctx, nil, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *gin.Context) error {
	params := &web.ChangeUserMobileRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)

	if uid == 2054 || uid == 2055 {
		return ginutil.BusinessError(ctx, "预览账号不支持修改手机号！")
	}

	if !u.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		return ginutil.BusinessError(ctx, "短信验证码填写错误！")
	}

	user, _ := u.service.Dao().FindById(jwtutil.GetUid(ctx))

	if user.Mobile != params.Mobile {
		return ginutil.BusinessError(ctx, "手机号与原手机号一致无需修改！")
	}

	if !encrypt.VerifyPassword(user.Password, params.Password) {
		return ginutil.BusinessError(ctx, "账号密码填写错误！")
	}

	_, err := u.service.Dao().BaseUpdate(&model.Users{}, entity.MapStrAny{"id": user.Id}, entity.MapStrAny{"mobile": params.Mobile})
	if err != nil {
		return ginutil.BusinessError(ctx, "手机号修改失败！")
	}

	return ginutil.Success(ctx, nil, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *gin.Context) error {
	params := &web.ChangeUserEmailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	// todo 1.验证邮件激活码是否正确

	return nil
}
