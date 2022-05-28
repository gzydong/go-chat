package organize

import (
	"go-chat/internal/dao"
	"go-chat/internal/model"
)

type IPositionDao interface {
	FindAll() ([]*model.OrganizePost, error)
}

type PositionDao struct {
	*dao.BaseDao
}

func NewPositionDao(baseDao *dao.BaseDao) *PositionDao {
	return &PositionDao{BaseDao: baseDao}
}

func (dao *PositionDao) FindAll() ([]*model.OrganizePost, error) {
	items := make([]*model.OrganizePost, 0)

	err := dao.Db().Model(model.OrganizePost{}).Where("status = 1").Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
