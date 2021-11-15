package service

import (
	"context"
	"errors"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"gorm.io/gorm"
	"time"
)

type ContactApplyService struct {
	*BaseService
}

func NewContactsApplyService(base *BaseService) *ContactApplyService {
	return &ContactApplyService{BaseService: base}
}

func (s *ContactApplyService) Create(ctx context.Context, uid int, req *request.ContactApplyCreateRequest) error {
	err := s.db.Create(model.UsersFriendsApply{
		UserId:    uid,
		FriendId:  req.FriendId,
		Remark:    req.Remarks,
		CreatedAt: time.Now(),
	}).Error

	return err
}

func (s *ContactApplyService) Accept(ctx context.Context, uid int, req *request.ContactApplyAcceptRequest) error {
	var err error

	var applyInfo model.UsersFriendsApply

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var userFriends model.UsersFriends
		if tx.Where("user_id = ? and friend_id = ?", applyInfo.UserId, applyInfo.FriendId).First(&userFriends).Error != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if tx.Create(&model.UsersFriends{
					UserId:   applyInfo.UserId,
					FriendId: applyInfo.FriendId,
					Remark:   req.Remarks,
					Status:   1,
				}).Error != nil {
					return err
				}
			}
		} else {
			tx.Model(&model.UsersFriends{}).Where("id = ?", userFriends.Id).Update("status", "1")
		}

		var userFriends2 model.UsersFriends
		if tx.Where("user_id = ? and friend_id = ?", applyInfo.FriendId, applyInfo.UserId).First(&userFriends2).Error != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if tx.Create(model.UsersFriends{
					UserId:   applyInfo.FriendId,
					FriendId: applyInfo.UserId,
					Remark:   applyInfo.Remark,
					Status:   1,
				}).Error != nil {
					return err
				}
			}
		} else {
			tx.Model(&model.UsersFriends{}).Where("id = ?", userFriends2.Id).Update("status", "1")
		}

		return tx.Delete(&model.UsersFriendsApply{}, applyInfo.Id).Error
	})

	return nil
}

// List 联系人申请列表
func (s *ContactApplyService) List(ctx context.Context, uid, page, size int) ([]*model.ApplyListItem, error) {
	fields := []string{
		"users_friends_apply.id",
		"users_friends_apply.remark",
		"users.nickname",
		"users.avatar",
		"users.mobile",
		"users_friends_apply.user_id",
		"users_friends_apply.friend_id",
		"users_friends_apply.created_at",
	}

	tx := s.db.Table("users_friends_apply")
	tx.Select(fields)
	tx.Joins("left join `users` ON `users`.id = users_friends_apply.user_id")
	tx.Where("users_friends_apply.friend_id = ?", uid)
	tx.Order("users_friends_apply.id desc")

	items := make([]*model.ApplyListItem, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
