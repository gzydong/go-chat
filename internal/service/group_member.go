package service

import (
	"go-chat/internal/entity"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type GroupMemberService struct {
	*BaseService
	repo *repo.GroupMember
}

func NewGroupMemberService(baseService *BaseService, repo *repo.GroupMember) *GroupMemberService {
	return &GroupMemberService{BaseService: baseService, repo: repo}
}

func (s *GroupMemberService) Dao() *repo.GroupMember {
	return s.repo
}

// ChangeGroupNickname 修改群名片
func (s *GroupMemberService) ChangeGroupNickname(groupId int, userId int, remark string) error {

	_, err := s.repo.BaseUpdate(&model.GroupMember{}, entity.MapStrAny{"group_id": groupId, "user_id": userId}, entity.MapStrAny{"user_card": remark})

	return err
}

// Handover 交接群主权限
func (s *GroupMemberService) Handover(groupId int, userId int, memberId int) error {
	return s.Db().Transaction(func(tx *gorm.DB) error {

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
	return s.Db().Model(model.GroupMember{}).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("leader", leader).Error
}

func (s *GroupMemberService) UpdateMuteStatus(groupId int, userId int, status int) error {
	return s.Db().Model(model.GroupMember{}).Where("group_id = ? and user_id = ?", groupId, userId).UpdateColumn("is_mute", status).Error
}
