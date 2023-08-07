package service

import (
	"context"
	"errors"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type ContactGroupService struct {
	*repo.Source
	contactGroup *repo.ContactGroup
}

func NewContactGroupService(source *repo.Source, contactGroup *repo.ContactGroup) *ContactGroupService {
	return &ContactGroupService{Source: source, contactGroup: contactGroup}
}

func (c *ContactGroupService) Delete(ctx context.Context, id int, uid int) error {
	return c.contactGroup.Txx(ctx, func(tx *gorm.DB) error {
		res := tx.Delete(&model.ContactGroup{}, "id = ? and user_id = ?", id, uid)
		if err := res.Error; err != nil {
			return err
		}

		if res.RowsAffected == 0 {
			return errors.New("数据不存在")
		}

		return tx.Table("contact").Where("user_id = ? and group_id = ?", uid, id).UpdateColumn("group_id", 0).Error
	})
}

// GetUserGroup 用户联系人分组列表
func (c *ContactGroupService) GetUserGroup(ctx context.Context, uid int) ([]*model.ContactGroup, error) {

	var items []*model.ContactGroup
	err := c.Db().WithContext(ctx).Table("contact_group").Where("user_id = ?", uid).Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
