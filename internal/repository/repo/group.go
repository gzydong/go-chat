package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Group struct {
	ichat.Repo[model.Group]
}

func NewGroup(db *gorm.DB) *Group {
	return &Group{Repo: ichat.NewRepo[model.Group](db)}
}

func (g *Group) SearchOvertList(ctx context.Context, name string, page, size int) ([]*model.Group, error) {
	return g.FindAll(ctx, func(db *gorm.DB) {
		if name != "" {
			db.Where("group_name LIKE ?", "%"+name+"%")
		} else {
			db.Where("is_overt = ?", 1)
		}

		db.Where("is_dismiss = 0").Order("created_at desc").Offset((page - 1) * size).Limit(size)
	})
}
