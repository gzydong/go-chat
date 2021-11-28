package service

import (
	"context"
	"go-chat/app/model"
)

type ContactService struct {
	*BaseService
}

func NewContactService(base *BaseService) *ContactService {
	return &ContactService{BaseService: base}
}

// EditRemark 编辑联系人备注
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) EditRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return s.db.Model(&model.UsersFriends{}).Where("user_id = ? and friend_id = ?", uid, friendId).Update("remark", remark).Error
}

// Delete 删除联系人
// @params uid      用户ID
// @params friendId 联系人ID
func (s *ContactService) Delete(ctx context.Context, uid, friendId int) error {
	return s.db.Model(model.UsersFriends{}).Where("user_id = ? and friend_id = ?", uid, friendId).Update("status", 0).Error
}

// List 获取联系人列表
// @params uid      用户ID
func (s *ContactService) List(ctx context.Context, uid int) ([]*model.ContactListItem, error) {

	items := make([]*model.ContactListItem, 0)

	tx := s.db.Table("users_friends")
	tx.Select([]string{
		"users.id",
		"users.nickname",
		"users.avatar",
		"users.motto",
		"users.gender",
		"users_friends.remark",
	})

	tx.Joins("inner join `users` ON `users`.id = users_friends.friend_id")
	tx.Where("`users_friends`.user_id = ? and users_friends.status = ?", uid, 1)

	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ContactService) GetContactIds(ctx context.Context, uid int) []int64 {
	var ids []int64

	s.db.Model(model.UsersFriends{}).Where("user_id = ? and status = ?", uid, 1).Pluck("friend_id", &ids)

	return ids
}
