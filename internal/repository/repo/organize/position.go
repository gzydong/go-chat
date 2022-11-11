package organize

import (
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type IPositionDao interface {
	FindAll() ([]*model.OrganizePost, error)
}

type Position struct {
	*repo.Base
}

func NewPosition(baseDao *repo.Base) *Position {
	return &Position{Base: baseDao}
}

func (repo *Position) FindAll() ([]*model.OrganizePost, error) {
	items := make([]*model.OrganizePost, 0)

	err := repo.Db().Model(model.OrganizePost{}).Where("status = 1").Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
