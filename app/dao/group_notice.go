package dao

import (
	"context"
	"go-chat/app/model"
	"gorm.io/gorm"
)

type GroupNoticeDao struct {
	*Base
}

func (dao *GroupNoticeDao) Db() *gorm.DB {
	return dao.db
}

func (dao *GroupNoticeDao) GetListAll(ctx context.Context, groupId int) ([]*model.SearchNoticeItem, error) {
	var items []*model.SearchNoticeItem

	fields := []string{
		"lar_group_notice.id",
		"lar_group_notice.creator_id",
		"lar_group_notice.title",
		"lar_group_notice.content",
		"lar_group_notice.is_top",
		"lar_group_notice.is_confirm",
		"lar_group_notice.confirm_users",
		"lar_group_notice.created_at",
		"lar_group_notice.updated_at",
		"lar_users.avatar",
		"lar_users.nickname",
	}

	err := dao.db.Table("lar_group_notice").
		Select(fields).
		Joins("left join lar_users on lar_users.id = lar_group_notice.creator_id").
		Where("lar_group_notice.group_id = ? and lar_group_notice.is_delete = ?", groupId, 0).
		Order("lar_group_notice.is_top desc").
		Order("lar_group_notice.updated_at desc").
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
