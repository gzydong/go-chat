package service

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type MemberItem struct {
	UserId   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
	Motto    string `json:"motto"`
	Leader   int    `json:"leader"`
	UserCard string `json:"user_card"`
}

type GroupMemberService struct {
	db *gorm.DB
}

func NewGroupMemberService(db *gorm.DB) *GroupMemberService {
	return &GroupMemberService{
		db: db,
	}
}

// isMember 判断用户是否是群成员
func (s *GroupMemberService) IsMember(groupId, userId int) bool {
	result := &model.GroupMember{}

	count := s.db.Select("id").Where("group_id = ? and user_id = ? and is_quit = ?", groupId, userId, 0).Unscoped().First(result).RowsAffected

	return count != 0
}

// GetMemberIds 获取所有群成员用户ID
func (s *GroupMemberService) GetMemberIds(groupId int) []int {
	var ids []int

	_ = s.db.Model(&model.GroupMember{}).Select("user_id").Where("group_id = ? and is_quit = ?", groupId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetMemberIds 获取所有群成员ID
func (s *GroupMemberService) GetUserGroupIds(userId int) []int {
	var ids []int

	_ = s.db.Model(&model.GroupMember{}).Select("id").Where("user_id = ? and is_quit = ?", userId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetGroupMembers 获取群组成员列表
func (s *GroupMemberService) GetGroupMembers(groupId int) []*MemberItem {
	var items []*MemberItem

	fields := []string{
		"group_member.leader",
		"group_member.user_card",
		"group_member.user_id",
		"users.avatar",
		"users.nickname",
		"users.gender",
		"users.motto",
	}

	s.db.Table("group_member").
		Select(fields).
		Joins("left join users on users.id = group_member.user_id").
		Where("group_member.group_id = ? and group_member.is_quit = ?", groupId, 0).
		Order("group_member.leader desc").
		Unscoped().
		Scan(&items)

	return items
}

// GetMemberRemarks 获取指定群成员的备注信息
func (s *GroupMemberService) GetMemberRemarks(groupId int, userId int) string {
	var remarks string

	s.db.Model(&model.GroupMember{}).
		Select("user_card").
		Where("group_id = ? and user_id = ?", groupId, userId).
		Unscoped().
		Scan(&remarks)

	return remarks
}

func (s *GroupMemberService) GetGroupMemberCount(gid int) int64 {
	num := int64(0)

	s.db.Model(&model.GroupMember{}).
		Where("group_id = ? and is_quit = ?", gid, 0).
		Unscoped().Count(&num)

	return num
}
