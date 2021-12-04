package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/timeutil"
)

type UserService struct {
	dao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{dao: userDao}
}

func (s *UserService) Dao() *dao.UserDao {
	return s.dao
}

// Register 注册用户
func (s *UserService) Register(param *request.RegisterRequest) (*model.User, error) {
	if s.dao.IsMobileExist(param.Mobile) {
		return nil, errors.New("账号已存在! ")
	}

	hash, _ := auth.Encrypt(param.Password)
	user, err := s.dao.Create(&model.User{
		Mobile:    param.Mobile,
		Nickname:  param.Nickname,
		Password:  hash,
		CreatedAt: timeutil.DateTime(),
		UpdatedAt: timeutil.DateTime(),
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 登录处理
func (s *UserService) Login(mobile string, password string) (*model.User, error) {
	user, err := s.dao.FindByMobile(mobile)
	if err != nil {
		return nil, errors.New("登录账号不存在! ")
	}

	if !auth.Compare(user.Password, password) {
		return nil, errors.New("登录密码填写错误! ")
	}

	return user, nil
}

// Forget 账号找回
func (s *UserService) Forget(input *request.ForgetRequest) (bool, error) {
	// 账号查询
	user, err := s.dao.FindByMobile(input.Mobile)
	if err != nil || user.ID == 0 {
		return false, errors.New("账号不存在! ")
	}

	// 生成 hash 密码
	hash, _ := auth.Encrypt(input.Password)

	_, err = s.Dao().BaseUpdate(&model.User{}, gin.H{"id": user.ID}, gin.H{"password": hash})
	if err != nil {
		return false, errors.New("密码修改失败！")
	}

	return true, nil
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {
	user := &model.User{}

	if ok, _ := s.dao.FindByIds(user, []int{uid}, "id,password"); !ok {
		return errors.New("用户不存在！")
	}

	if !auth.Compare(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	hash, _ := auth.Encrypt(password)

	_, err := s.dao.BaseUpdate(&model.User{}, gin.H{"id": user.ID}, gin.H{"password": hash})
	if err != nil {
		return err
	}

	return nil
}
