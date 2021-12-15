package dao

import (
	"go-chat/app/model"
)

type UserDao struct {
	*BaseDao
}

func NewUserDao(baseDao *BaseDao) *UserDao {
	return &UserDao{BaseDao: baseDao}
}

// Create 创建数据
func (dao *UserDao) Create(user *model.Users) (*model.Users, error) {
	if err := dao.Db().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindById ID查询
func (dao *UserDao) FindById(userid int) (*model.Users, error) {
	user := &model.Users{}

	if err := dao.Db().Where(&model.Users{Id: userid}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (dao *UserDao) FindByMobile(mobile string) (*model.Users, error) {
	user := &model.Users{}

	if err := dao.Db().Where(&model.Users{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// IsMobileExist 判断手机号是否存在
func (dao *UserDao) IsMobileExist(mobile string) bool {
	user := &model.Users{}

	rowsAffects := dao.Db().Select("id").Where(&model.Users{Mobile: mobile}).First(user).RowsAffected

	return rowsAffects != 0
}
