package service

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IUserService = (*UserService)(nil)

type IUserService interface {
	Register(ctx context.Context, opt *UserRegisterOpt) (*model.Users, error)
	Login(ctx context.Context, mobile string, password string) (*model.Users, error)
	Forget(ctx context.Context, opt *UserForgetOpt) (bool, error)
	UpdatePassword(ctx context.Context, uid int, oldPassword string, password string) error
	OauthBind(ctx context.Context, mobile string, oauthUser *model.OAuthUser) (int, error)
}

type UserService struct {
	UsersRepo      *repo.Users
	OAuthUsersRepo *repo.OAuthUsers
}

type UserRegisterOpt struct {
	Nickname string
	Mobile   string
	Password string
	Platform string
}

// Register 注册用户
func (s *UserService) Register(ctx context.Context, opt *UserRegisterOpt) (*model.Users, error) {
	if s.UsersRepo.IsMobileExist(ctx, opt.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	user := &model.Users{
		Mobile:    lo.ToPtr(opt.Mobile),
		Nickname:  opt.Nickname,
		Avatar:    "",
		Gender:    model.UsersGenderDefault,
		Password:  encrypt.HashPassword(opt.Password),
		Motto:     "",
		Email:     "",
		Birthday:  "",
		IsRobot:   model.No,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.UsersRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(ctx context.Context, mobile string, password string) (*model.Users, error) {
	user, err := s.UsersRepo.FindByMobile(ctx, mobile)
	if err != nil {
		if utils.IsSqlNoRows(err) {
			return nil, entity.ErrAccountOrPassword
		}

		return nil, err
	}

	if !encrypt.VerifyPassword(user.Password, password) {
		return nil, entity.ErrAccountOrPassword
	}

	if user.IsDisabled() {
		return nil, entity.ErrAccountDisabled
	}

	return user, nil
}

// UserForgetOpt ForgetRequest 账号找回接口验证
type UserForgetOpt struct {
	Mobile   string
	Password string
	SmsCode  string
}

// Forget 账号找回
func (s *UserService) Forget(ctx context.Context, opt *UserForgetOpt) (bool, error) {
	user, err := s.UsersRepo.FindByMobile(ctx, opt.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	affected, err := s.UsersRepo.UpdateById(context.TODO(), user.Id, map[string]any{
		"password": encrypt.HashPassword(opt.Password),
	})

	return affected > 0, err
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(ctx context.Context, uid int, oldPassword string, password string) error {
	user, err := s.UsersRepo.FindById(ctx, uid)
	if err != nil {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	_, err = s.UsersRepo.UpdateById(ctx, user.Id, map[string]any{
		"password": encrypt.HashPassword(password),
	})

	return err
}

func (s *UserService) OauthBind(ctx context.Context, mobile string, oauthUser *model.OAuthUser) (int, error) {
	userinfo, err := s.UsersRepo.FindByMobile(ctx, mobile)
	if err != nil && !utils.IsSqlNoRows(err) {
		return 0, err
	}

	if userinfo != nil {
		oauth, err := s.OAuthUsersRepo.FindByWhere(ctx, "user_id = ? and oauth_type = ?", userinfo.Id, oauthUser.OAuthType)
		if err != nil && !utils.IsSqlNoRows(err) {
			return 0, err
		}

		if oauth != nil && oauth.UserId == int32(userinfo.Id) {
			return userinfo.Id, nil
		}

		if oauth != nil {
			return 0, entity.ErrAccountBinded
		}

		_, err = s.OAuthUsersRepo.UpdateByWhere(ctx, map[string]any{
			"user_id": userinfo.Id,
		}, "id = ?", oauthUser.Id)
		if err != nil {
			return 0, err
		}

		return userinfo.Id, nil
	}

	user := &model.Users{
		Mobile:    lo.ToPtr(mobile),
		Nickname:  oauthUser.Nickname,
		Avatar:    oauthUser.Avatar,
		Gender:    3,
		IsRobot:   1,
		Status:    model.UsersStatusNormal,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err = s.UsersRepo.Txx(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.OAuthUser{}).Where("id = ?", oauthUser.Id).Updates(&model.OAuthUser{
			UserId: int32(user.Id),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return user.Id, nil
}
