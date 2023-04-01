package service

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type ContactService struct {
	*repo.Source
	contact *repo.Contact
}

func NewContactService(source *repo.Source, contact *repo.Contact) *ContactService {
	return &ContactService{Source: source, contact: contact}
}

func (s *ContactService) Dao() *repo.Contact {
	return s.contact
}

// EditRemark 编辑联系人备注
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) EditRemark(ctx context.Context, uid int, friendId int, remark string) error {

	_, err := s.contact.UpdateWhere(ctx, map[string]any{"remark": remark}, "user_id = ? and friend_id = ?", uid, friendId)
	if err == nil {
		_ = s.contact.SetFriendRemark(ctx, uid, friendId, remark)
	}

	return err
}

// Delete 删除联系人
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) Delete(ctx context.Context, uid, friendId int) error {
	return s.contact.Model(ctx).Where("user_id = ? and friend_id = ?", uid, friendId).Update("status", 0).Error
}

// List 获取联系人列表
// @params uid      用户ID
func (s *ContactService) List(ctx context.Context, uid int) ([]*model.ContactListItem, error) {

	tx := s.contact.Model(ctx)
	tx.Select([]string{
		"users.id",
		"users.nickname",
		"users.avatar",
		"users.motto",
		"users.gender",
		"contact.remark",
		"contact.group_id",
	})
	tx.Joins("inner join `users` ON `users`.id = contact.friend_id")
	tx.Where("contact.user_id = ? and contact.status = ?", uid, 1)

	var items []*model.ContactListItem
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ContactService) GetContactIds(ctx context.Context, uid int) []int64 {

	var ids []int64
	s.contact.Model(ctx).Where("user_id = ? and status = ?", uid, 1).Pluck("friend_id", &ids)

	return ids
}

func (s *ContactService) MoveGroup(ctx context.Context, uid int, friendId int, groupId int) error {
	contact, err := s.Dao().FindByWhere(ctx, "user_id = ? and friend_id  = ?", uid, friendId)
	if err != nil {
		return err
	}

	return s.Db().Transaction(func(tx *gorm.DB) error {
		if contact.GroupId > 0 {
			err := tx.Table("contact_group").Where("id = ? and user_id = ?", contact.GroupId, uid).Updates(map[string]any{
				"num": gorm.Expr("num - 1"),
			}).Error

			if err != nil {
				return err
			}
		}

		err := tx.Table("contact").Where("user_id = ? and friend_id = ? and group_id = ?", uid, friendId, contact.GroupId).UpdateColumn("group_id", groupId).Error
		if err != nil {
			return err
		}

		return tx.Table("contact_group").Where("id = ? and user_id = ?", groupId, uid).Updates(map[string]any{
			"num": gorm.Expr("num + 1"),
		}).Error
	})
}
