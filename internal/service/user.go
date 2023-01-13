package service

import (
	"context"
	"errors"

	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repo.Users
}

func NewUserService(repo *repo.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Dao() *repo.Users {
	return s.repo
}

type UserRegisterOpt struct {
	Nickname string
	Mobile   string
	Password string
	Platform string
}

// Register 注册用户
func (s *UserService) Register(opts *UserRegisterOpt) (*model.Users, error) {
	if s.repo.IsMobileExist(opts.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	user, err := s.repo.Create(&model.Users{
		Mobile:   opts.Mobile,
		Nickname: opts.Nickname,
		Password: encrypt.HashPassword(opts.Password),
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(mobile string, password string) (*model.Users, error) {
	user, err := s.repo.FindByMobile(mobile)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("登录账号不存在! ")
		}

		return nil, err
	}

	if !encrypt.VerifyPassword(user.Password, password) {
		return nil, errors.New("登录密码填写错误! ")
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
func (s *UserService) Forget(opts *UserForgetOpt) (bool, error) {

	user, err := s.repo.FindByMobile(opts.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	err = s.Dao().Db.Model(&model.Users{}).
		Where("id = ?", user.Id).
		Update("password", encrypt.HashPassword(opts.Password)).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {

	user, err := s.Dao().FindById(context.Background(), uid)
	if err != nil {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	err = s.Dao().Db.Model(&model.Users{}).Where("id = ?", user.Id).Update("password", encrypt.HashPassword(password)).Error
	if err != nil {
		return err
	}

	return nil
}
