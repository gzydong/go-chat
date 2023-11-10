package service

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IGroupMemberService = (*GroupMemberService)(nil)

type IGroupMemberService interface {
	Handover(ctx context.Context, groupId int, userId int, memberId int) error
	SetLeaderStatus(ctx context.Context, groupId int, userId int, leader int) error
	SetMuteStatus(ctx context.Context, groupId int, userId int, status int) error
}

type GroupMemberService struct {
	*repo.Source
	GroupMemberRepo *repo.GroupMember
}

func (g *GroupMemberService) Handover(ctx context.Context, groupId int, userId int, memberId int) error {
	return g.Source.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&model.GroupMember{}).Where("group_id = ? and user_id = ? and leader = 2", groupId, userId).Update("leader", 0).Error
		if err != nil {
			return err
		}

		err = tx.Model(&model.GroupMember{}).Where("group_id = ? and user_id = ?", groupId, memberId).Update("leader", 2).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (g *GroupMemberService) SetLeaderStatus(ctx context.Context, groupId int, userId int, leader int) error {
	return g.GroupMemberRepo.Model(ctx).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("leader", leader).Error
}

func (g *GroupMemberService) SetMuteStatus(ctx context.Context, groupId int, userId int, status int) error {
	return g.GroupMemberRepo.Model(ctx).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("is_mute", status).Error
}
