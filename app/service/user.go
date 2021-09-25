package service

import (
	"errors"
	"fmt"
	"go-chat/app/helper"
	"go-chat/app/model"
	"go-chat/app/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

// Register 注册用户
func (s *UserService) Register(username string, password string) (bool, error) {
	hash, err := helper.GeneratePassword([]byte(password))
	if err != nil {
		return false, err
	}

	fmt.Println(hash)

	return true, nil
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
