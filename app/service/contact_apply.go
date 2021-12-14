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
	err := s.db.Create(&model.UsersFriendsApply{
		UserId:    uid,
		FriendId:  req.FriendId,
		Remark:    req.Remarks,
		CreatedAt: time.Now(),
	}).Error

	return err
}

// Accept 同意好友申请
func (s *ContactApplyService) Accept(ctx context.Context, uid int, req *request.ContactApplyAcceptRequest) error {
	var (
		err       error
		applyInfo *model.UsersFriendsApply
	)

	if err := s.db.First(&applyInfo, req.ApplyId).Error; err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		addFriendFunc := func(uid, fid int, remark string) error {
			var friends *model.UsersFriends

			err = tx.Where("user_id = ? and friend_id = ?", uid, fid).First(&friends).Error

			// 数据存在则更新
			if err == nil {
				return tx.Model(&model.UsersFriends{}).Where("id = ?", friends.Id).Updates(&model.UsersFriends{
					Remark: remark,
					Status: 1,
				}).Error
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			return tx.Create(&model.UsersFriends{
				UserId:   uid,
				FriendId: fid,
				Remark:   remark,
				Status:   1,
			}).Error
		}

		var user *model.Users
		if err := tx.Select("id", "nickname").First(&user, applyInfo.FriendId).Error; err != nil {
			return err
		}

		if err := addFriendFunc(applyInfo.UserId, applyInfo.FriendId, user.Nickname); err != nil {
			return err
		}

		if err := addFriendFunc(applyInfo.FriendId, applyInfo.UserId, req.Remarks); err != nil {
			return err
		}

		return tx.Delete(&model.UsersFriendsApply{}, applyInfo.Id).Error
	})

	return nil
}

// Decline 拒绝好友申请
func (s *ContactApplyService) Decline(ctx context.Context, uid int, req *request.ContactApplyDeclineRequest) error {
	return s.db.Delete(&model.UsersFriendsApply{}, "id = ? and friend_id = ?", req.ApplyId, uid).Error
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
