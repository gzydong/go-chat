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
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register 注册用户
func (s *UserService) Register(param *request.RegisterRequest) (*model.User, error) {
	exist := s.repo.IsMobileExist(param.Mobile)
	if exist {
		return nil, errors.New("账号已存在! ")
	}

	hash, err := helper.GeneratePassword([]byte(param.Password))
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(&model.User{
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
func (s *UserService) Login(mobile string, password string) (*model.User, error) {
	user, err := s.repo.FindByMobile(mobile)
	if err != nil {
		return nil, errors.New("登录账号不存在! ")
	}

	if !helper.VerifyPassword([]byte(password), []byte(user.Password)) {
		return nil, errors.New("登录密码填写错误! ")
	}

	return user, nil
}

// Forget 账号找回
func (s *UserService) Forget(input *request.ForgetRequest) (bool, error) {
	// 账号查询
	user, err := s.repo.FindByMobile(input.Mobile)
	if err != nil || user.ID == 0 {
		return false, errors.New("账号不存在! ")
	}

	// 生成 hash 密码
	hash, _ := helper.GeneratePassword([]byte(input.Password))

	_, err = s.repo.Update(&model.User{
		ID: user.ID,
	}, map[string]interface{}{
		"password": hash,
	})

	if err != nil {
		return false, errors.New("密码修改失败！")
	}

	return true, nil
}
