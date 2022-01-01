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

// Dao 获取数据 dao 层
func (s *GroupMemberService) Dao() *dao.GroupMemberDao {
	return s.dao
}

// nolint UpdateMemberCard 修改群名片
func (s *GroupMemberService) UpdateMemberCard(groupId int, userId int, remark string) error {

	_, err := s.dao.BaseUpdate(&model.GroupMember{}, entity.Map{"group_id": groupId, "user_id": userId}, entity.Map{"user_card": remark})
	if err == nil {
		// todo 更新缓存
	}

	return err
}
