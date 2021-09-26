package service

import (
	"errors"
	"go-chat/app/http/request"
	"time"

	"go-chat/app/helper"
	"go-chat/app/model"
	"go-chat/app/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

// Register 注册用户
func (s *UserService) Register(param *request.RegisterRequest) (*model.User, error) {
	exist := s.Repo.IsMobileExist(param.Mobile)
	if exist {
		return nil, errors.New("账号已存在! ")
	}

	// todo 这里需要判断短信验证码是否正确

	hash, err := helper.GeneratePassword([]byte(param.Password))
	if err != nil {
		return nil, err
	}

	user, err := s.Repo.Create(&model.User{
		Mobile:    param.Mobile,
		Nickname:  param.Nickname,
		Password:  string(hash),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(username string, password string) (*model.User, error) {
	user, err := s.Repo.FindByMobile(username)
	if err != nil {
		return nil, errors.New("登录账号不存在! ")
	}

	if !helper.VerifyPassword([]byte(password), []byte(user.Password)) {
		return nil, errors.New("登录密码填写错误! ")
	}

	return user, nil
}
