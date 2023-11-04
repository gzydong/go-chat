package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type GroupNotice struct {
	ichat.Repo[model.GroupNotice]
}

func NewGroupNotice(db *gorm.DB) *GroupNotice {
	return &GroupNotice{Repo: ichat.NewRepo[model.GroupNotice](db)}
}

func (g *GroupNotice) GetListAll(ctx context.Context, groupId int) ([]*model.SearchNoticeItem, error) {

	fields := []string{
		"group_notice.id",
		"group_notice.creator_id",
		"group_notice.title",
		"group_notice.content",
		"group_notice.is_top",
		"group_notice.is_confirm",
		"group_notice.confirm_users",
		"group_notice.created_at",
		"group_notice.updated_at",
		"users.avatar",
		"users.nickname",
	}

	query := g.Repo.Db.WithContext(ctx).Table("group_notice")
	query.Joins("left join users on users.id = group_notice.creator_id")
	query.Where("group_notice.group_id = ? and group_notice.is_delete = ?", groupId, 0)
	query.Order("group_notice.is_top desc")
	query.Order("group_notice.created_at desc")

	var items []*model.SearchNoticeItem
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// GetLatestNotice 获取最新公告
func (g *GroupNotice) GetLatestNotice(ctx context.Context, groupId int) (*model.GroupNotice, error) {
	var info model.GroupNotice
	err := g.Repo.Db.WithContext(ctx).Last(&info, "group_id = ? and is_delete = ?", groupId, 0).Error
	if err != nil {
		return nil, err
	}

	return &info, nil
}
