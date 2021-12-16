package dao

import "go-chat/app/model"

type GroupDao struct {
	*BaseDao
}

func NewGroupDao(baseDao *BaseDao) *GroupDao {
	return &GroupDao{BaseDao: baseDao}
}

func (dao *GroupDao) FindById(id int) (*model.Group, error) {
	info := &model.Group{}

	if err := dao.Db().First(&info, id).Error; err != nil {
		return nil, err
	}

	return info, nil
}
