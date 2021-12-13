package service

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/model"
)

type GroupMemberService struct {
	dao *dao.GroupMemberDao
}

func NewGroupMemberService(dao *dao.GroupMemberDao) *GroupMemberService {
	return &GroupMemberService{dao: dao}
}

// Dao 获取数据 dao 层
func (s *GroupMemberService) Dao() *dao.GroupMemberDao {
	return s.dao
}

// nolint UpdateMemberCard 修改群名片
func (s *GroupMemberService) UpdateMemberCard(groupId int, userId int, remark string) error {

	_, err := s.dao.BaseUpdate(&model.GroupMember{}, gin.H{"group_id": groupId, "user_id": userId}, gin.H{"user_card": remark})
	if err == nil {
		// todo 更新缓存
	}

	return err
}
