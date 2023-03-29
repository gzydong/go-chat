package service

import (
	"context"
	"errors"
	"fmt"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
)

type ContactApplyService struct {
	*repo.Source
}

func NewContactApplyService(source *repo.Source) *ContactApplyService {
	return &ContactApplyService{Source: source}
}

type ContactApplyCreateOpt struct {
	UserId   int
	Remarks  string
	FriendId int
}

func (s *ContactApplyService) Create(ctx context.Context, opt *ContactApplyCreateOpt) error {

	apply := &model.ContactApply{
		UserId:   opt.UserId,
		FriendId: opt.FriendId,
		Remark:   opt.Remarks,
	}

	if err := s.Db().Create(apply).Error; err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.EventContactApply,
		"data": jsonutil.Encode(map[string]any{
			"apply_id": int64(apply.Id),
			"type":     1,
		}),
	}

	s.Redis().Incr(ctx, fmt.Sprintf("friend-apply:user_%d", opt.FriendId))

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))

	return nil
}

type ContactApplyAcceptOpt struct {
	UserId  int
	Remarks string
	ApplyId int
}

// Accept 同意好友申请
func (s *ContactApplyService) Accept(ctx context.Context, opt *ContactApplyAcceptOpt) (*model.ContactApply, error) {
	var (
		err       error
		applyInfo *model.ContactApply
	)

	if err := s.Db().First(&applyInfo, "id = ? and friend_id = ?", opt.ApplyId, opt.UserId).Error; err != nil {
		return nil, err
	}

	err = s.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		addFriendFunc := func(uid, fid int, remark string) error {
			var friends *model.Contact

			err = tx.Where("user_id = ? and friend_id = ?", uid, fid).First(&friends).Error

			// 数据存在则更新
			if err == nil {
				return tx.Model(&model.Contact{}).Where("id = ?", friends.Id).Updates(&model.Contact{
					Remark: remark,
					Status: 1,
				}).Error
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			return tx.Create(&model.Contact{
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

		if err := addFriendFunc(applyInfo.FriendId, applyInfo.UserId, opt.Remarks); err != nil {
			return err
		}

		return tx.Delete(&model.ContactApply{}, "user_id = ? and friend_id = ?", applyInfo.UserId, applyInfo.FriendId).Error
	})

	return applyInfo, err
}

type ContactApplyDeclineOpt struct {
	UserId  int
	Remarks string
	ApplyId int
}

// Decline 拒绝好友申请
func (s *ContactApplyService) Decline(ctx context.Context, opt *ContactApplyDeclineOpt) error {
	err := s.Db().Delete(&model.ContactApply{}, "id = ? and friend_id = ?", opt.ApplyId, opt.UserId).Error
	if err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.EventContactApply,
		"data": jsonutil.Encode(map[string]any{
			"apply_id": int64(opt.ApplyId),
			"type":     2,
		}),
	}

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))
	return nil
}

// List 联系人申请列表
func (s *ContactApplyService) List(ctx context.Context, uid int) ([]*model.ApplyItem, error) {
	fields := []string{
		"contact_apply.id",
		"contact_apply.remark",
		"users.nickname",
		"users.avatar",
		"users.mobile",
		"contact_apply.user_id",
		"contact_apply.friend_id",
		"contact_apply.created_at",
	}

	tx := s.Db().WithContext(ctx).Table("contact_apply")
	tx.Joins("left join `users` ON `users`.id = contact_apply.user_id")
	tx.Where("contact_apply.friend_id = ?", uid)
	tx.Order("contact_apply.id desc")

	var items []*model.ApplyItem
	if err := tx.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ContactApplyService) GetApplyUnreadNum(ctx context.Context, uid int) int {

	num, err := s.Redis().Get(ctx, fmt.Sprintf("friend-apply:user_%d", uid)).Int()
	if err != nil {
		return 0
	}

	return num
}

func (s *ContactApplyService) ClearApplyUnreadNum(ctx context.Context, uid int) {
	s.Redis().Del(ctx, fmt.Sprintf("friend-apply:user_%d", uid))
}
