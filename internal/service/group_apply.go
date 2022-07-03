package service

import (
	"context"
	"errors"

	"go-chat/internal/repository/dao"
	"go-chat/internal/repository/model"
)

type GroupApplyService struct {
	*BaseService
	dao *dao.GroupApplyDao
}

func NewGroupApplyService(baseService *BaseService, dao *dao.GroupApplyDao) *GroupApplyService {
	return &GroupApplyService{BaseService: baseService, dao: dao}
}

func (s *GroupApplyService) Dao() *dao.GroupApplyDao {
	return s.dao
}

func (s *GroupApplyService) Auth(ctx context.Context, applyId, userId int) bool {
	info := &model.GroupApply{}

	err := s.Db().First(info, "id = ?", applyId).Error
	if err != nil {
		return false
	}

	member := &model.GroupMember{}
	err = s.Db().First(member, "group_id = ? and user_id = ? and leader in (1,2) and is_quit = 0", info.GroupId).Error
	if err != nil {
		return false
	}

	return member.Id == 0
}

func (s *GroupApplyService) Insert(ctx context.Context, groupId, userId int, remark string) error {
	return s.Db().Create(&model.GroupApply{
		GroupId: groupId,
		UserId:  userId,
		Remark:  remark,
	}).Error
}

func (s *GroupApplyService) Delete(ctx context.Context, applyId, userId int) error {

	if !s.Auth(ctx, applyId, userId) {
		return errors.New("auth failed")
	}

	return s.Db().Delete(&model.GroupApply{}, "id = ?", applyId).Error
}
