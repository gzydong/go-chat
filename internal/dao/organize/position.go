package organize

import (
	"go-chat/internal/dao"
	"go-chat/internal/model"
)

type PositionDao struct {
	*dao.BaseDao
}

type PositionDaoInterface interface {
	FindAll() ([]*model.OrganizePost, error)
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
