package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Group struct {
	core.Repo[model.Group]
}

func NewGroup(db *gorm.DB) *Group {
	return &Group{Repo: core.NewRepo[model.Group](db)}
}

type SearchOvertListOpt struct {
	Name   string
	UserId int
	Page   int
	Size   int
}

func (g *Group) SearchOvertList(ctx context.Context, opt *SearchOvertListOpt) ([]*model.Group, error) {
	return g.Repo.FindAll(ctx, func(db *gorm.DB) {
		if opt.Name != "" {
			db.Where("name like ?", "%"+opt.Name+"%")
		}

		db.Where("is_overt = ?", model.Yes)
		db.Where("id NOT IN (?)", g.Repo.Db.Select("group_id").Where("user_id = ? and is_quit = ?", opt.UserId, model.No).Table("group_member"))
		db.Where("is_dismiss = ?", model.No).Order("created_at desc").Offset((opt.Page - 1) * opt.Size).Limit(opt.Size)
	})
}
