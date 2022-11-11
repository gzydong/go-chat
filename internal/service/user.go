package service

import (
	"errors"

	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type UserService struct {
	dao *repo.Users
}

func NewUserService(userDao *repo.Users) *UserService {
	return &UserService{dao: userDao}
}

func (s *UserService) Dao() *repo.Users {
	return s.dao
}

type UserRegisterOpt struct {
	Nickname string
	Mobile   string
	Password string
	Platform string
}

// Register 注册用户
func (s *UserService) Register(opts *UserRegisterOpt) (*model.Users, error) {
	if s.dao.IsMobileExist(opts.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	hash, _ := encrypt.HashPassword(opts.Password)
	user, err := s.dao.Create(&model.Users{
		Mobile:   opts.Mobile,
		Nickname: opts.Nickname,
		Password: hash,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(mobile string, password string) (*model.Users, error) {
	user, err := s.dao.FindByMobile(mobile)
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

	user, err := s.dao.FindByMobile(opts.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	// 生成 hash 密码
	hash, _ := encrypt.HashPassword(opts.Password)

	err = s.Dao().Db().Model(&model.Users{}).Where("id = ?", user.Id).Update("password", hash).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {

	user, err := s.Dao().FindById(uid)
	if err != nil {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	hash, err := encrypt.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.Dao().Db().Model(&model.Users{}).Where("id = ?", user.Id).Update("password", hash).Error
	if err != nil {
		return err
	}

	return nil
}
