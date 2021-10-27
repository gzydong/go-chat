package dao

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

// Create 创建数据
func (u *UserDao) Create(user *model.User) (*model.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Update 更新数据
func (u *UserDao) Update(user *model.User, column map[string]interface{}) (int64, error) {
	res := u.db.Model(&user).Updates(column)

	if err := res.Error; err != nil {
		return int64(0), err
	}

	return res.RowsAffected, nil
}

// FindById ID查询
func (u *UserDao) FindById(userid int) (*model.User, error) {
	user := &model.User{}
	if err := u.db.Where(&model.User{ID: userid}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (u *UserDao) FindByMobile(mobile string) (*model.User, error) {
	user := &model.User{}
	if err := u.db.Where(&model.User{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// IsMobileExist 判断手机号是否存在
func (u *UserDao) IsMobileExist(mobile string) bool {
	user := &model.User{}

	rowsAffects := u.db.Select("id").Where(&model.User{Mobile: mobile}).First(user).RowsAffected

	return rowsAffects != 0
}
