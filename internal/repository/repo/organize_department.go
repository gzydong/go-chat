package repo

import (
	"context"

	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Department struct {
	core.Repo[model.OrganizeDept]
}

func NewDepartment(db *gorm.DB) *Department {
	return &Department{Repo: core.NewRepo[model.OrganizeDept](db)}
}

func (d *Department) List(ctx context.Context) ([]*model.OrganizeDept, error) {
	return d.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("is_deleted = 1").Order("parent_id asc,order_num asc")
	})
}
