package organize

import (
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type IPosition interface {
	FindAll() ([]*model.OrganizePost, error)
}

type Position struct {
	*repo.Base
}

func NewPosition(base *repo.Base) *Position {
	return &Position{Base: base}
}

func (repo *Position) FindAll() ([]*model.OrganizePost, error) {
	items := make([]*model.OrganizePost, 0)

	err := repo.Db.Model(model.OrganizePost{}).Where("status = 1").Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
