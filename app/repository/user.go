package repository

import (
	"go-chat/app/model"
	"go-chat/connect"
)

type UserRepository struct {
	DB *connect.MySQL
}

// FindByMobile 手机号查询
func (u *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	user := &model.User{}
	if err := u.DB.Db.Where(&model.User{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
