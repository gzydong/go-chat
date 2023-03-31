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
	users *repo.Users
}

func NewUserService(repo *repo.Users) *UserService {
	return &UserService{users: repo}
}

func (s *UserService) Dao() *repo.Users {
	return s.users
}

type UserRegisterOpt struct {
	Nickname string
	Mobile   string
	Password string
	Platform string
}

// Register 注册用户
func (s *UserService) Register(opts *UserRegisterOpt) (*model.Users, error) {
	if s.users.IsMobileExist(opts.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	return s.users.Create(&model.Users{
		Mobile:   opts.Mobile,
		Nickname: opts.Nickname,
		Password: encrypt.HashPassword(opts.Password),
	})
}

// Login 登录处理
func (s *UserService) Login(mobile string, password string) (*model.Users, error) {

	user, err := s.users.FindByMobile(mobile)
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

	user, err := s.users.FindByMobile(opts.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	affected, err := s.users.UpdateById(context.TODO(), user.Id, map[string]any{
		"password": encrypt.HashPassword(opts.Password),
	})

	return affected > 0, err
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {

	user, err := s.users.FindById(context.TODO(), uid)
	if err != nil {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	_, err = s.users.UpdateById(context.TODO(), user.Id, map[string]any{
		"password": encrypt.HashPassword(password),
	})

	return err
}
