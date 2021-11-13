package dao

import (
	"context"
	"go-chat/app/model"
)

type GroupNoticeDao struct {
	*Base
}

func (dao *GroupNoticeDao) GetListAll(ctx context.Context, groupId int) ([]*model.SearchNoticeItem, error) {
	var items []*model.SearchNoticeItem

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

	err := dao.Db.Table("group_notice").
		Select(fields).
		Joins("left join users on users.id = group_notice.creator_id").
		Where("group_notice.group_id = ? and group_notice.is_delete = ?", groupId, 0).
		Order("group_notice.is_top desc").
		Order("group_notice.updated_at desc").
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
