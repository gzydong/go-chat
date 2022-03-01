package service

import (
	"errors"

	"go-chat/internal/dao"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/encrypt"
)

type UserRegisterOpts struct {
	Nickname string
	Mobile   string
	Password string
	SmsCode  string
	Platform string
}

// ForgetRequest 账号找回接口验证
type UserForgetOpts struct {
	Mobile   string
	Password string
	SmsCode  string
}

type UserService struct {
	dao *dao.UsersDao
}

func NewUserService(userDao *dao.UsersDao) *UserService {
	return &UserService{dao: userDao}
}

func (s *UserService) Dao() *dao.UsersDao {
	return s.dao
}

// Register 注册用户
func (s *UserService) Register(opts *UserRegisterOpts) (*model.Users, error) {
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
		return nil, errors.New("登录账号不存在! ")
	}

	if !encrypt.VerifyPassword(user.Password, password) {
		return nil, errors.New("登录密码填写错误! ")
	}

	return user, nil
}

// Forget 账号找回
func (s *UserService) Forget(opts *UserForgetOpts) (bool, error) {

	user, err := s.dao.FindByMobile(opts.Mobile)
	if err != nil || user.Id == 0 {
		return false, errors.New("账号不存在! ")
	}

	// 生成 hash 密码
	hash, _ := encrypt.HashPassword(opts.Password)

	_, err = s.Dao().BaseUpdate(&model.Users{}, entity.Map{"id": user.Id}, entity.Map{"password": hash})
	if err != nil {
		return false, errors.New("密码修改失败！")
	}

	return true, nil
}

// UpdatePassword 修改用户密码
func (s *UserService) UpdatePassword(uid int, oldPassword string, password string) error {
	user := &model.Users{}

	if ok, _ := s.dao.FindByIds(user, []int{uid}, "id,password"); !ok {
		return errors.New("用户不存在！")
	}

	if !encrypt.VerifyPassword(user.Password, oldPassword) {
		return errors.New("密码验证不正确！")
	}

	hash, _ := encrypt.HashPassword(password)

	_, err := s.dao.BaseUpdate(&model.Users{}, entity.Map{"id": user.Id}, entity.Map{"password": hash})
	if err != nil {
		return err
	}

	return nil
}
