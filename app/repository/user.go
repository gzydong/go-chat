package repository

import (
	"go-chat/app/model"
	"go-chat/connect"
)

type UserRepository struct {
	DB *connect.MySQL
}

// findByMobile 手机号查询
func (u *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	var user model.User

	result := u.DB.Db.Where(&model.User{Mobile: mobile}).First(&user)

	return &user, result.Error
}
