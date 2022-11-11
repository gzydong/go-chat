package repo

import (
	"context"

	"go-chat/internal/repository/model"
)

type GroupNotice struct {
	*Base
}

func NewGroupNotice(baseDao *Base) *GroupNotice {
	return &GroupNotice{Base: baseDao}
}

func (repo *GroupNotice) GetListAll(ctx context.Context, groupId int) ([]*model.SearchNoticeItem, error) {

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

	query := repo.Db().Table("group_notice")
	query.Joins("left join users on users.id = group_notice.creator_id")
	query.Where("group_notice.group_id = ? and group_notice.is_delete = ?", groupId, 0)
	query.Order("group_notice.is_top desc")
	query.Order("group_notice.created_at desc")

	items := make([]*model.SearchNoticeItem, 0)
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// GetLatestNotice 获取最新公告
func (repo *GroupNotice) GetLatestNotice(ctx context.Context, groupId int) (*model.GroupNotice, error) {
	info := &model.GroupNotice{}

	err := repo.Db().Last(info, "group_id = ? and is_delete = ?", groupId, 0).Error
	if err != nil {
		return nil, err
	}

	return info, nil
}
