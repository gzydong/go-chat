package v1

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

var _ web.IUserHandler = (*User)(nil)

type User struct {
	Redis        *redis.Client
	UsersRepo    *repo.Users
	OrganizeRepo *repo.Organize
	UserService  service.IUserService
	SmsService   service.ISmsService
	Rsa          rsautil.IRsa
}

func (u *User) Detail(ctx context.Context, _ *web.UserDetailRequest) (*web.UserDetailResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	user, err := u.UsersRepo.FindByIdWithCache(ctx, int(session.UserId))
	if err != nil {
		return nil, err
	}

	return &web.UserDetailResponse{
		Mobile:   lo.FromPtr(user.Mobile),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		Birthday: user.Birthday,
	}, nil
}

func (u *User) Setting(ctx context.Context, req *web.UserSettingRequest) (*web.UserSettingResponse, error) {
	session, err := middleware.FormContext[entity.WebClaims](ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.UsersRepo.FindByIdWithCache(ctx, int(session.UserId))
	if err != nil {
		return nil, err
	}

	isOk, err := u.OrganizeRepo.IsQiyeMember(ctx, int(session.UserId))
	if err != nil {
		return nil, err
	}

	return &web.UserSettingResponse{
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
	}, nil
}

func (u *User) DetailUpdate(ctx context.Context, req *web.UserDetailUpdateRequest) (*web.UserDetailUpdateResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	if req.Birthday != "" {
		if !timeutil.IsDate(req.Birthday) {
			return nil, errorx.New(400, "birthday 错误")
		}
	}

	uid := session.UserId

	_, err := u.UsersRepo.UpdateById(ctx, uid, map[string]any{
		"nickname": strings.TrimSpace(strings.ReplaceAll(req.Nickname, " ", "")),
		"avatar":   req.Avatar,
		"gender":   req.Gender,
		"motto":    req.Motto,
		"birthday": req.Birthday,
	})

	if err != nil {
		return nil, err
	}

	_ = u.UsersRepo.ClearTableCache(ctx, int(uid))
	return &web.UserDetailUpdateResponse{}, nil
}

func (u *User) PasswordUpdate(ctx context.Context, in *web.UserPasswordUpdateRequest) (*web.UserPasswordUpdateResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	uid := session.UserId
	if uid == 2054 || uid == 2055 {
		return nil, entity.ErrPermissionDenied
	}

	oldPassword, err := u.Rsa.Decrypt(in.OldPassword)
	if err != nil {
		return nil, err
	}

	newPassword, err := u.Rsa.Decrypt(in.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := u.UserService.UpdatePassword(ctx, int(uid), string(oldPassword), string(newPassword)); err != nil {
		return nil, err
	}

	_ = u.UsersRepo.ClearTableCache(ctx, int(uid))
	return nil, nil
}

func (u *User) MobileUpdate(ctx context.Context, in *web.UserMobileUpdateRequest) (*web.UserMobileUpdateResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	uid := session.UserId

	user, _ := u.UsersRepo.FindById(ctx, uid)
	if lo.FromPtr(user.Mobile) == in.Mobile {
		return nil, errorx.New(400, "手机号与原手机号一致无需修改")
	}

	password, err := u.Rsa.Decrypt(in.Password)
	if err != nil {
		return nil, err
	}

	if !encrypt.VerifyPassword(user.Password, string(password)) {
		return nil, entity.ErrAccountOrPasswordError
	}

	if uid == 2054 || uid == 2055 {
		return nil, entity.ErrPermissionDenied
	}

	if !u.SmsService.Verify(ctx, entity.SmsChangeAccountChannel, in.Mobile, in.SmsCode) {
		return nil, entity.ErrSmsCodeError
	}

	_, err = u.UsersRepo.UpdateById(ctx, user.Id, map[string]any{
		"mobile": in.Mobile,
	})

	if err != nil {
		return nil, err
	}

	_ = u.UsersRepo.ClearTableCache(ctx, user.Id)
	return nil, nil
}

func (u *User) EmailUpdate(ctx context.Context, req *web.UserEmailUpdateRequest) (*web.UserEmailUpdateResponse, error) {
	//TODO implement me
	panic("implement me")
}
