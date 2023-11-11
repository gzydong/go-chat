package v1

import (
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type User struct {
	UsersRepo    *repo.Users
	OrganizeRepo *repo.Organize

	UserService service.IUserService
	SmsService  service.ISmsService
}

// Detail 个人用户信息
func (u *User) Detail(ctx *ichat.Context) error {

	user, err := u.UsersRepo.FindById(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&web.UserDetailResponse{
		Id:       int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		Birthday: user.Birthday,
	})
}

// Setting 用户设置
func (u *User) Setting(ctx *ichat.Context) error {

	uid := ctx.UserId()

	user, err := u.UsersRepo.FindById(ctx.Ctx(), uid)
	if err != nil {
		return ctx.Error(err.Error())
	}

	isOk, err := u.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid)
	if err != nil {
		return ctx.Error(err.Error())
	}

	return ctx.Success(&web.UserSettingResponse{
		UserInfo: &web.UserSettingResponse_UserInfo{
			Uid:      int32(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Motto:    user.Motto,
			Gender:   int32(user.Gender),
			IsQiye:   isOk,
			Mobile:   user.Mobile,
			Email:    user.Email,
		},
		Setting: &web.UserSettingResponse_ConfigInfo{},
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *ichat.Context) error {

	params := &web.UserDetailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.Birthday != "" {
		if !timeutil.IsDateFormat(params.Birthday) {
			return ctx.InvalidParams("birthday 格式错误")
		}
	}

	_, err := u.UsersRepo.UpdateById(ctx.Ctx(), ctx.UserId(), map[string]any{
		"nickname": strings.TrimSpace(strings.Replace(params.Nickname, " ", "", -1)),
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Motto,
		"birthday": params.Birthday,
	})

	if err != nil {
		return ctx.ErrorBusiness("个人信息修改失败！")
	}

	return ctx.Success(nil, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *ichat.Context) error {

	params := &web.UserPasswordUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if uid == 2054 || uid == 2055 {
		return ctx.ErrorBusiness("预览账号不支持修改密码！")
	}

	if err := u.UserService.UpdatePassword(ctx.UserId(), params.OldPassword, params.NewPassword); err != nil {
		return ctx.ErrorBusiness("密码修改失败！")
	}

	return ctx.Success(nil, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *ichat.Context) error {

	params := &web.UserMobileUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	user, _ := u.UsersRepo.FindById(ctx.Ctx(), uid)
	if user.Mobile == params.Mobile {
		return ctx.ErrorBusiness("手机号与原手机号一致无需修改！")
	}

	if !encrypt.VerifyPassword(user.Password, params.Password) {
		return ctx.ErrorBusiness("账号密码填写错误！")
	}

	if uid == 2054 || uid == 2055 {
		return ctx.ErrorBusiness("预览账号不支持修改手机号！")
	}

	if !u.SmsService.Verify(ctx.Ctx(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		return ctx.ErrorBusiness("短信验证码填写错误！")
	}

	_, err := u.UsersRepo.UpdateById(ctx.Ctx(), user.Id, map[string]any{
		"mobile": params.Mobile,
	})

	if err != nil {
		return ctx.ErrorBusiness("手机号修改失败！")
	}

	return ctx.Success(nil, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *ichat.Context) error {

	params := &web.UserEmailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// todo 1.验证邮件激活码是否正确

	return ctx.ErrorBusiness("手机号修改成功！")
}
