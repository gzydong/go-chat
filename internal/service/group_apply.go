package service

import (
	"context"
	"errors"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

var _ IGroupApplyService = (*GroupApplyService)(nil)

type IGroupApplyService interface {
	Auth(ctx context.Context, applyId, userId int) bool
	Insert(ctx context.Context, groupId, userId int, remark string) error
	Delete(ctx context.Context, applyId, userId int) error
}

type GroupApplyService struct {
	*repo.Source
	GroupApplyRepo *repo.GroupApply
}

func (s *GroupApplyService) Auth(ctx context.Context, applyId, userId int) bool {
	info, err := s.GroupApplyRepo.FindById(ctx, applyId)
	if err != nil {
		return false
	}

	var member model.GroupMember
	err = s.Source.Db().Debug().WithContext(ctx).Select("id").First(&member, "group_id = ? and user_id = ? and leader in (1,2) and is_quit = 0", info.GroupId, userId).Error

	return err == nil && member.Id > 0
}

func (s *GroupApplyService) Insert(ctx context.Context, groupId, userId int, remark string) error {
	return s.GroupApplyRepo.Create(ctx, &model.GroupApply{
		GroupId: groupId,
		UserId:  userId,
		Remark:  remark,
	})
}

func (s *GroupApplyService) Delete(ctx context.Context, applyId, userId int) error {

	if !s.Auth(ctx, applyId, userId) {
		return errors.New("auth failed")
	}

	return s.Source.Db().WithContext(ctx).Delete(&model.GroupApply{}, "id = ?", applyId).Error
}
