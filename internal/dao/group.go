package dao

import (
	"context"

	"go-chat/internal/model"
)

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

func (dao *GroupDao) SearchOvertList(ctx context.Context, name string, page, size int) ([]*model.Group, error) {

	tx := dao.Db().Table("group")

	if name != "" {
		tx.Where("group_name LIKE ?", "%"+name+"%")
	} else {
		tx.Where("is_overt = ?", 1)
	}

	items := make([]*model.Group, 0)
	err := tx.Where("is_dismiss = 0").Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
