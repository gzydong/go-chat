package service

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type ContactService struct {
	*BaseService
	repo *repo.Contact
}

func NewContactService(baseService *BaseService, dao *repo.Contact) *ContactService {
	return &ContactService{BaseService: baseService, repo: dao}
}

func (s *ContactService) Dao() *repo.Contact {
	return s.repo
}

// EditRemark 编辑联系人备注
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) EditRemark(ctx context.Context, uid int, friendId int, remark string) error {

	err := s.repo.Model(ctx).Where("user_id = ? and friend_id = ?", uid, friendId).Update("remark", remark).Error
	if err == nil {
		_ = s.repo.SetFriendRemark(ctx, uid, friendId, remark)
	}

	return err
}

// Delete 删除联系人
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) Delete(ctx context.Context, uid, friendId int) error {
	return s.repo.Model(ctx).Where("user_id = ? and friend_id = ?", uid, friendId).Update("status", 0).Error
}

// List 获取联系人列表
// @params uid      用户ID
func (s *ContactService) List(ctx context.Context, uid int) ([]*model.ContactListItem, error) {

	tx := s.repo.Model(ctx)

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

	items := make([]*model.ContactListItem, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ContactService) GetContactIds(ctx context.Context, uid int) []int64 {

	var ids []int64
	s.repo.Model(ctx).Where("user_id = ? and status = ?", uid, 1).Pluck("friend_id", &ids)

	return ids
}
