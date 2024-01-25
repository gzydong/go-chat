package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Position struct {
	core.Repo[model.OrganizePost]
}

func NewPosition(db *gorm.DB) *Position {
	return &Position{Repo: core.NewRepo[model.OrganizePost](db)}
}

func (p *Position) List(ctx context.Context) ([]*model.OrganizePost, error) {
	return p.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("status = 1").Order("sort asc")
	})
}
