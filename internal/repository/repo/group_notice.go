package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type GroupNotice struct {
	core.Repo[model.GroupNotice]
}

func NewGroupNotice(db *gorm.DB) *GroupNotice {
	return &GroupNotice{Repo: core.NewRepo[model.GroupNotice](db)}
}

// GetLatestNotice 获取最新公告
func (g *GroupNotice) GetLatestNotice(ctx context.Context, groupId int) (*model.GroupNotice, error) {
	var info model.GroupNotice
	err := g.Repo.Db.WithContext(ctx).Last(&info, "group_id = ?", groupId).Error
	if err != nil {
		return nil, err
	}

	return &info, nil
}
