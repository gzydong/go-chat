package dao

import (
	"context"

	"go-chat/internal/model"
)

type GroupApplyDao struct {
	*BaseDao
}

func NewGroupApply(baseDao *BaseDao) *GroupApplyDao {
	return &GroupApplyDao{BaseDao: baseDao}
}

func (dao *GroupApplyDao) List(ctx context.Context, groupId int) ([]*model.GroupApplyList, error) {

	fields := []string{
		"group_apply.id",
		"group_apply.group_id",
		"group_apply.user_id",
		"group_apply.remark",
		"group_apply.created_at",
		"users.avatar",
		"users.nickname",
	}

	query := dao.Db().Table("group_apply")
	query.Joins("left join users on users.id = group_apply.user_id")
	query.Where("group_apply.group_id = ?", groupId)
	query.Order("group_apply.created_at desc")

	items := make([]*model.GroupApplyList, 0)
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
