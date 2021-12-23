package service

import (
	"context"
	"errors"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/jsonutil"
	"gorm.io/gorm"
)

type ContactApplyService struct {
	*BaseService
}

func NewContactsApplyService(base *BaseService) *ContactApplyService {
	return &ContactApplyService{BaseService: base}
}

func (s *ContactApplyService) Create(ctx context.Context, uid int, req *request.ContactApplyCreateRequest) error {

	apply := &model.ContactApply{
		UserId:   uid,
		FriendId: req.FriendId,
		Remark:   req.Remarks,
	}

	if err := s.db.Create(apply).Error; err != nil {
		return err
	}

	body := map[string]interface{}{
		"event": entity.EventFriendApply,
		"data": jsonutil.JsonEncode(map[string]interface{}{
			"apply_id": int64(apply.Id),
			"type":     1,
		}),
	}

	s.rds.Publish(ctx, entity.IMGatewayAll, jsonutil.JsonEncode(body))

	return nil
}

// Accept 同意好友申请
func (s *ContactApplyService) Accept(ctx context.Context, uid int, req *request.ContactApplyAcceptRequest) error {
	var (
		err       error
		applyInfo *model.ContactApply
	)

	if err := s.db.First(&applyInfo, req.ApplyId).Error; err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
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

		if err := addFriendFunc(applyInfo.FriendId, applyInfo.UserId, req.Remarks); err != nil {
			return err
		}

		return tx.Delete(&model.ContactApply{}, applyInfo.Id).Error
	})

	if err == nil {
		body := map[string]interface{}{
			"event": entity.EventFriendApply,
			"data": jsonutil.JsonEncode(map[string]interface{}{
				"apply_id": int64(applyInfo.Id),
				"type":     2,
			}),
		}

		s.rds.Publish(ctx, entity.IMGatewayAll, jsonutil.JsonEncode(body))
	}

	return err
}

// Decline 拒绝好友申请
func (s *ContactApplyService) Decline(ctx context.Context, uid int, req *request.ContactApplyDeclineRequest) error {
	err := s.db.Delete(&model.ContactApply{}, "id = ? and friend_id = ?", req.ApplyId, uid).Error

	if err == nil {
		body := map[string]interface{}{
			"event": entity.EventFriendApply,
			"data": jsonutil.JsonEncode(map[string]interface{}{
				"apply_id": int64(req.ApplyId),
				"type":     2,
			}),
		}

		s.rds.Publish(ctx, entity.IMGatewayAll, jsonutil.JsonEncode(body))
	}

	return err
}

// List 联系人申请列表
func (s *ContactApplyService) List(ctx context.Context, uid, page, size int) ([]*model.ApplyItem, error) {
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

	tx := s.db.Table("contact_apply")
	tx.Select(fields)
	tx.Joins("left join `users` ON `users`.id = contact_apply.user_id")
	tx.Where("contact_apply.friend_id = ?", uid)
	tx.Order("contact_apply.id desc")

	items := make([]*model.ApplyItem, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
