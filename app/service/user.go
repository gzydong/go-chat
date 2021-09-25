package service

import (
	"errors"
	"fmt"

	"go-chat/app/helper"
	"go-chat/app/model"
)

type UserService struct {
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
	hashPassword := "tea123jas"

	// ...数据库查询

	if helper.VerifyPassword([]byte(password), []byte(hashPassword)) {
		return nil, errors.New("登录密码填写错误")
	}

	return &model.User{}, nil
}
