package v1

import (
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type User struct {
	Redis        *redis.Client
	UsersRepo    *repo.Users
	OrganizeRepo *repo.Organize
	UserService  service.IUserService
	SmsService   service.ISmsService
	Rsa          rsautil.IRsa
}

// Detail 个人用户信息
func (u *User) Detail(ctx *core.Context) error {
	user, err := u.UsersRepo.FindByIdWithCache(ctx.GetContext(), ctx.GetAuthId())
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.UserDetailResponse{
		Mobile:   lo.FromPtr(user.Mobile),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		Birthday: user.Birthday,
	})
}

// Setting 用户设置
func (u *User) Setting(ctx *core.Context) error {

	uid := ctx.GetAuthId()

	user, err := u.UsersRepo.FindByIdWithCache(ctx.GetContext(), uid)
	if err != nil {
		return ctx.Error(err)
	}

	isOk, err := u.OrganizeRepo.IsQiyeMember(ctx.GetContext(), uid)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.UserSettingResponse{
		UserInfo: &web.UserSettingResponse_UserInfo{
			Uid:      int32(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Motto:    user.Motto,
			Gender:   int32(user.Gender),
			IsQiye:   isOk,
			Mobile:   lo.FromPtr(user.Mobile),
			Email:    user.Email,
		},
		Setting: &web.UserSettingResponse_ConfigInfo{},
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *core.Context) error {
	in := &web.UserDetailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.Birthday != "" {
		if !timeutil.IsDate(in.Birthday) {
			return ctx.InvalidParams("birthday 格式错误")
		}
	}

	uid := ctx.GetAuthId()
	_, err := u.UsersRepo.UpdateById(ctx.GetContext(), ctx.GetAuthId(), map[string]any{
		"nickname": strings.TrimSpace(strings.Replace(in.Nickname, " ", "", -1)),
		"avatar":   in.Avatar,
		"gender":   in.Gender,
		"motto":    in.Motto,
		"birthday": in.Birthday,
	})

	if err != nil {
		return ctx.Error(err)
	}

	_ = u.UsersRepo.ClearTableCache(ctx.GetContext(), uid)

	return ctx.Success(nil, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *core.Context) error {
	in := &web.UserPasswordUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.GetAuthId()
	if uid == 2054 || uid == 2055 {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	oldPassword, err := u.Rsa.Decrypt(in.OldPassword)
	if err != nil {
		return ctx.Error(err)
	}

	newPassword, err := u.Rsa.Decrypt(in.NewPassword)
	if err != nil {
		return ctx.Error(err)
	}

	if err := u.UserService.UpdatePassword(ctx.GetContext(), ctx.GetAuthId(), string(oldPassword), string(newPassword)); err != nil {
		return ctx.Error(err)
	}

	_ = u.UsersRepo.ClearTableCache(ctx.GetContext(), uid)

	return ctx.Success(nil, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *core.Context) error {
	in := &web.UserMobileUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.GetAuthId()

	user, _ := u.UsersRepo.FindById(ctx.GetContext(), uid)
	if lo.FromPtr(user.Mobile) == in.Mobile {
		return ctx.InvalidParams("手机号与原手机号一致无需修改！")
	}

	password, err := u.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if !encrypt.VerifyPassword(user.Password, string(password)) {
		return ctx.Error(entity.ErrAccountOrPasswordError)
	}

	if uid == 2054 || uid == 2055 {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	if !u.SmsService.Verify(ctx.GetContext(), entity.SmsChangeAccountChannel, in.Mobile, in.SmsCode) {
		return ctx.Error(entity.ErrSmsCodeError)
	}

	_, err = u.UsersRepo.UpdateById(ctx.GetContext(), user.Id, map[string]any{
		"mobile": in.Mobile,
	})

	if err != nil {
		return ctx.Error(err)
	}

	_ = u.UsersRepo.ClearTableCache(ctx.GetContext(), user.Id)

	return ctx.Success(nil, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *core.Context) error {
	in := &web.UserEmailUpdateRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	// todo 1.验证邮件激活码是否正确

	return ctx.Success(nil, "手机号修改成功！")
}
