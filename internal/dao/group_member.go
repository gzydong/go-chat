package dao

import (
	"context"

	"go-chat/internal/cache"
	"go-chat/internal/model"
)

type GroupMemberDao struct {
	*BaseDao
	relation *cache.Relation
}

func NewGroupMemberDao(baseDao *BaseDao, relation *cache.Relation) *GroupMemberDao {
	return &GroupMemberDao{BaseDao: baseDao, relation: relation}
}

// IsMaster 判断是否是群主
func (dao *GroupMemberDao) IsMaster(gid, uid int) bool {
	result := &model.GroupMember{}

	count := dao.Db().Select("id").Where("group_id = ? and user_id = ? and leader = 2 and is_quit = 0", gid, uid).First(result).RowsAffected

	return count == 1
}

// IsLeader 判断是否是群主或管理员
func (dao *GroupMemberDao) IsLeader(gid, uid int) bool {
	result := &model.GroupMember{}

	count := dao.Db().Select("id").Where("group_id = ? and user_id = ? and leader in (1,2) and is_quit = 0", gid, uid).First(result).RowsAffected

	return count == 1
}

// IsMember 检测是属于群成员
func (dao *GroupMemberDao) IsMember(gid, uid int, cache bool) bool {
	if cache && dao.relation.IsGroupRelation(context.Background(), uid, gid) == nil {
		return true
	}

	result := &model.GroupMember{}

	count := dao.Db().Select("id").Where("group_id = ? and user_id = ? and is_quit = ?", gid, uid, 0).First(result).RowsAffected

	if count == 1 {
		dao.relation.SetGroupRelation(context.Background(), uid, gid)
	}

	return count != 0
}

// GetMemberIds 获取所有群成员用户ID
func (dao *GroupMemberDao) GetMemberIds(groupId int) []int {
	ids := make([]int, 0)

	_ = dao.Db().Model(&model.GroupMember{}).Select("user_id").Where("group_id = ? and is_quit = ?", groupId, 0).Scan(&ids)

	return ids
}

// GetUserGroupIds 获取所有群成员ID
func (dao *GroupMemberDao) GetUserGroupIds(uid int) []int {
	ids := make([]int, 0)

	_ = dao.Db().Model(&model.GroupMember{}).Where("user_id = ? and is_quit = ?", uid, 0).Pluck("group_id", &ids)

	return ids
}

// CountMemberTotal 统计群成员总数
func (dao *GroupMemberDao) CountMemberTotal(gid int) int64 {
	num := int64(0)

	dao.Db().Model(&model.GroupMember{}).Where("group_id = ? and is_quit = ?", gid, 0).Count(&num)

	return num
}

// GetMemberRemark 获取指定群成员的备注信息
func (dao *GroupMemberDao) GetMemberRemark(groupId int, userId int) string {
	var remarks string

	dao.Db().Model(&model.GroupMember{}).Select("user_card").Where("group_id = ? and user_id = ?", groupId, userId).Scan(&remarks)

	return remarks
}

// GetMembers 获取群组成员列表
func (dao *GroupMemberDao) GetMembers(groupId int) []*model.MemberItem {
	fields := []string{
		"group_member.id",
		"group_member.leader",
		"group_member.user_card",
		"group_member.user_id",
		"group_member.is_mute",
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

type CountGroupMember struct {
	GroupId int `gorm:"column:group_id;"`
	Count   int `gorm:"column:count;"`
}

func (dao *GroupMemberDao) CountGroupMemberNum(ids []int) ([]*CountGroupMember, error) {
	items := make([]*CountGroupMember, 0)

	err := dao.Db().Table("group_member").Select("group_id,count(*) as count").Where("group_id in ? and is_quit = 0", ids).Group("group_id").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (dao *GroupMemberDao) CheckUserGroup(ids []int, userId int) ([]int, error) {
	items := make([]int, 0)

	err := dao.Db().Table("group_member").Select("group_id").Where("group_id in ? and user_id = ? and is_quit = 0", ids, userId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
