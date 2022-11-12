package repo

import (
	"fmt"

	"go-chat/internal/repository/model"
)

type Users struct {
	*Base
}

func NewUsers(base *Base) *Users {
	return &Users{Base: base}
}

// Create 创建数据
func (u *Users) Create(user *model.Users) (*model.Users, error) {
	if err := u.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindById ID查询
func (u *Users) FindById(uid int) (*model.Users, error) {

	if uid == 0 {
		return nil, fmt.Errorf("uid is empty")
	}

	user := &model.Users{}

	if err := u.Db.Where(&model.Users{Id: uid}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (u *Users) FindByMobile(mobile string) (*model.Users, error) {

	if len(mobile) == 0 {
		return nil, fmt.Errorf("mobile is empty")
	}

	user := &model.Users{}

	if err := u.Db.Where(&model.Users{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// IsMobileExist 判断手机号是否存在
func (u *Users) IsMobileExist(mobile string) bool {

	if len(mobile) == 0 {
		return false
	}

	user := &model.Users{}

	rowsAffects := u.Db.Select("id").Where(&model.Users{Mobile: mobile}).First(user).RowsAffected

	return rowsAffects != 0
}
