package service

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type GroupMemberService struct {
	db *gorm.DB
}

// isMember 判断用户是否是群成员
func (s *GroupMemberService) isMember(groupId, userId int) bool {
	result := &model.GroupMember{}

	count := s.db.Select("id").
		Where("group_id = ? and user_id = ? and is_quit = ?", groupId, userId, 0).Unscoped().
		First(result).RowsAffected

	return count != 0
}

// GetMemberIds 获取所有群成员ID
func (s *GroupMemberService) GetMemberIds(groupId int) []int {
	var ids []int

	_ = s.db.Model(&model.GroupMember{}).Select("id").Where("group_id = ? and is_quit = ?", groupId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetMemberIds 获取所有群成员ID
func (s *GroupMemberService) GetUserGroupIds(userId int) []int {
	var ids []int

	_ = s.db.Debug().Model(&model.GroupMember{}).Select("id").Where("user_id = ? and is_quit = ?", userId, 0).Unscoped().Scan(&ids)

	return ids
}
