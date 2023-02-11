package service

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type ContactGroupService struct {
	*BaseService
	repo *repo.ContactGroup
}

func NewContactGroupService(baseService *BaseService, repo *repo.ContactGroup) *ContactGroupService {
	return &ContactGroupService{BaseService: baseService, repo: repo}
}

func (c *ContactGroupService) Repo() *repo.ContactGroup {
	return c.repo
}

func (c *ContactGroupService) Delete(ctx context.Context, id int, uid int) error {
	return c.repo.Txx(ctx, func(tx *gorm.DB) error {
		res := tx.Delete(&model.ContactGroup{}, "id = ? and user_id = ?", id, uid)
		if err := res.Error; err != nil {
			return err
		}

		res = tx.Table(model.Contact{}.TableName()).Where("user_id = ? and group_id = ?", uid, id).UpdateColumn("group_id", 0)
		if err := res.Error; err != nil {
			return err
		}

		return nil
	})
}

func (c *ContactGroupService) Sort(ctx context.Context, uid int, values []*model.ContactGroup) error {
	return c.repo.Txx(ctx, func(tx *gorm.DB) error {
		for _, value := range values {
			err := tx.Table(model.ContactGroup{}.TableName()).Where("id = ? and user_id = ?", value.Id, uid).
				UpdateColumn("sort", value.Sort).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// GetUserGroup 用户联系人分组列表
func (c *ContactGroupService) GetUserGroup(ctx context.Context, uid int) ([]*model.ContactGroup, error) {

	var items []*model.ContactGroup

	err := c.db.WithContext(ctx).Table("contact_group").Where("user_id = ?", uid).Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
