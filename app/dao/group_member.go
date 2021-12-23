package dao

import (
	"context"
	"go-chat/app/cache"
	"go-chat/app/model"
)

type GroupMemberDao struct {
	*BaseDao
	relation *cache.Relation
}

func NewGroupMemberDao(baseDao *BaseDao, relation *cache.Relation) *GroupMemberDao {
	return &GroupMemberDao{BaseDao: baseDao, relation: relation}
}

// IsMember 检测是属于群成员
func (dao *GroupMemberDao) IsMember(gid, uid int, cache bool) bool {
	if dao.relation.IsGroupRelation(context.Background(), uid, gid) == nil {
		return true
	}

	result := &model.GroupMember{}

	count := dao.Db().Select("id").Where("group_id = ? and user_id = ? and is_quit = ?", gid, uid, 0).Unscoped().First(result).RowsAffected

	if count == 1 {
		dao.relation.SetGroupRelation(context.Background(), uid, gid)
	}

	return count != 0
}

// GetMemberIds 获取所有群成员用户ID
func (dao *GroupMemberDao) GetMemberIds(groupId int) []int {
	ids := make([]int, 0)

	_ = dao.Db().Model(&model.GroupMember{}).Select("user_id").Where("group_id = ? and is_quit = ?", groupId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetUserGroupIds 获取所有群成员ID
func (dao *GroupMemberDao) GetUserGroupIds(uid int) []int {
	ids := make([]int, 0)

	_ = dao.Db().Model(&model.GroupMember{}).Where("user_id = ? and is_quit = ?", uid, 0).Unscoped().Pluck("group_id", &ids)

	return ids
}

// CountMemberTotal 统计群成员总数
func (dao *GroupMemberDao) CountMemberTotal(gid int) int64 {
	num := int64(0)

	dao.Db().Model(&model.GroupMember{}).Where("group_id = ? and is_quit = ?", gid, 0).Unscoped().Count(&num)

	return num
}

// GetMemberRemark 获取指定群成员的备注信息
func (dao *GroupMemberDao) GetMemberRemark(groupId int, userId int) string {
	var remarks string

	dao.Db().Model(&model.GroupMember{}).Select("user_card").Where("group_id = ? and user_id = ?", groupId, userId).Unscoped().Scan(&remarks)

	return remarks
}

// GetMembers 获取群组成员列表
func (dao *GroupMemberDao) GetMembers(groupId int) []*model.MemberItem {
	fields := []string{
		"group_member.leader",
		"group_member.user_card",
		"group_member.user_id",
		"users.avatar",
		"users.nickname",
		"users.gender",
		"users.motto",
	}

	tx := dao.Db().Table("group_member")
	tx.Joins("left join users on users.id = group_member.user_id")
	tx.Where("group_member.group_id = ? and group_member.is_quit = ?", groupId, 0)
	tx.Order("group_member.leader desc")

	items := make([]*model.MemberItem, 0)
	tx.Unscoped().Select(fields).Scan(&items)

	return items
}
