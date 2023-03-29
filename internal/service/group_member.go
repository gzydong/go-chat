package service

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type GroupMemberService struct {
	*repo.Source
	member *repo.GroupMember
}

func NewGroupMemberService(source *repo.Source, repo *repo.GroupMember) *GroupMemberService {
	return &GroupMemberService{Source: source, member: repo}
}

func (s *GroupMemberService) Dao() *repo.GroupMember {
	return s.member
}

// Handover 交接群主权限
func (s *GroupMemberService) Handover(ctx context.Context, groupId int, userId int, memberId int) error {
	return s.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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

func (s *GroupMemberService) UpdateLeaderStatus(groupId int, userId int, leader int) error {
	return s.member.Model(context.TODO()).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("leader", leader).Error
}

func (s *GroupMemberService) UpdateMuteStatus(groupId int, userId int, status int) error {
	return s.member.Model(context.TODO()).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("is_mute", status).Error
}
