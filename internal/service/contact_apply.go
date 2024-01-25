package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
)

var _ IContactApplyService = (*ContactApplyService)(nil)

type IContactApplyService interface {
	Create(ctx context.Context, opt *ContactApplyCreateOpt) error
	Accept(ctx context.Context, opt *ContactApplyAcceptOpt) (*model.ContactApply, error)
	Decline(ctx context.Context, opt *ContactApplyDeclineOpt) error
	List(ctx context.Context, uid int) ([]*model.ApplyItem, error)
	GetApplyUnreadNum(ctx context.Context, uid int) int
	ClearApplyUnreadNum(ctx context.Context, uid int)
}

type ContactApplyService struct {
	*repo.Source
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

	if err := s.Source.Db().WithContext(ctx).Create(apply).Error; err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.SubEventContactApply,
		"data": jsonutil.Encode(map[string]any{
			"apply_id": apply.Id,
			"type":     1,
		}),
	}

	_, _ = s.Source.Redis().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, fmt.Sprintf("im:contact:apply:%d", opt.FriendId))
		pipe.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))
		return nil
	})

	return nil
}

type ContactApplyAcceptOpt struct {
	UserId  int
	Remarks string
	ApplyId int
}

// Accept 同意好友申请
func (s *ContactApplyService) Accept(ctx context.Context, opt *ContactApplyAcceptOpt) (*model.ContactApply, error) {

	db := s.Source.Db().WithContext(ctx)

	var applyInfo model.ContactApply
	if err := db.First(&applyInfo, "id = ? and friend_id = ?", opt.ApplyId, opt.UserId).Error; err != nil {
		return nil, err
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		addFriendFunc := func(uid, fid int, remark string) error {
			var contact model.Contact
			err := tx.Where("user_id = ? and friend_id = ?", uid, fid).First(&contact).Error

			// 数据存在则更新
			if err == nil {
				return tx.Model(&model.Contact{}).Where("id = ?", contact.Id).Updates(&model.Contact{
					Remark: remark,
					Status: 1,
				}).Error
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			return tx.Create(&model.Contact{
				UserId:    uid,
				FriendId:  fid,
				Remark:    remark,
				Status:    1,
				GroupId:   0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}).Error
		}

		var user model.Users
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

	return &applyInfo, err
}

type ContactApplyDeclineOpt struct {
	UserId  int
	Remarks string
	ApplyId int
}

// Decline 拒绝好友申请
func (s *ContactApplyService) Decline(ctx context.Context, opt *ContactApplyDeclineOpt) error {
	err := s.Source.Db().WithContext(ctx).Delete(&model.ContactApply{}, "id = ? and friend_id = ?", opt.ApplyId, opt.UserId).Error
	if err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.SubEventContactApply,
		"data": jsonutil.Encode(map[string]any{
			"apply_id": opt.ApplyId,
			"type":     2,
		}),
	}

	s.Source.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))
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

	tx := s.Source.Db().WithContext(ctx).Table("contact_apply")
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

	num, err := s.Source.Redis().Get(ctx, fmt.Sprintf("im:contact:apply:%d", uid)).Int()
	if err != nil {
		return 0
	}

	return num
}

func (s *ContactApplyService) ClearApplyUnreadNum(ctx context.Context, uid int) {
	s.Source.Redis().Del(ctx, fmt.Sprintf("im:contact:apply:%d", uid))
}
