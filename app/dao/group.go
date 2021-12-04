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

	dao.Db.First(&info, id)

	return info, nil
}
