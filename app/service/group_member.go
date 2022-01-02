package service

import (
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/model"
)

type GroupMemberService struct {
	*BaseService
	dao *dao.GroupMemberDao
}

func NewGroupMemberService(baseService *BaseService, dao *dao.GroupMemberDao) *GroupMemberService {
	return &GroupMemberService{BaseService: baseService, dao: dao}
}

func (s *GroupMemberService) Dao() *dao.GroupMemberDao {
	return s.dao
}

// EditMemberCard 修改群名片
func (s *GroupMemberService) EditMemberCard(groupId int, userId int, remark string) error {

	_, err := s.dao.BaseUpdate(&model.GroupMember{}, entity.Map{"group_id": groupId, "user_id": userId}, entity.Map{"user_card": remark})

	return err
}
