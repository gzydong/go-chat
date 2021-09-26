package repository

import (
	"go-chat/app/model"
	"go-chat/connect"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *connect.MySQL
}

// Create 创建数据
func (u *UserRepository) Create(user *model.User) (*model.User, error) {
	if err := u.db().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (u *UserRepository) FindByMobile(mobile string) (*model.User, error) {
	user := &model.User{}
	if err := u.db().Where(&model.User{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// IsMobileExist 判断手机号是否存在
func (u *UserRepository) IsMobileExist(mobile string) bool {
	user := &model.User{}

	rowsAffects := u.db().Select("id", "mobile").Where(&model.User{Mobile: mobile}).First(user).RowsAffected

	return rowsAffects != 0
}

func (u *UserRepository) db() *gorm.DB {
	return u.DB.Db
}
