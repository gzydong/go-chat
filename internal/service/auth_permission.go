package service

import (
	"context"

	"go-chat/internal/dao"
	"go-chat/internal/dao/organize"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/logger"
)

type AuthPermissionService struct {
	contactDao     *dao.ContactDao
	groupMemberDao *dao.GroupMemberDao
	organizeDao    *organize.OrganizeDao
}

func NewAuthPermissionService(contactDao *dao.ContactDao, groupMemberDao *dao.GroupMemberDao, organizeDao *organize.OrganizeDao) *AuthPermissionService {
	return &AuthPermissionService{contactDao: contactDao, groupMemberDao: groupMemberDao, organizeDao: organizeDao}
}

type AuthPermission struct {
	TalkType   int
	UserId     int
	ReceiverId int
}

func (a *AuthPermissionService) IsAuth(ctx context.Context, prem *AuthPermission) bool {
	if prem.TalkType == entity.ChatPrivateMode {
		// 这里需要判断双方是否都是企业成员，如果是则无需添加好友即可聊天
		if isOk, err := a.organizeDao.IsQiyeMember(prem.UserId, prem.ReceiverId); err != nil {
			logger.Error("[AuthPermission IsAuth] 查询数据异常 err: ", err)
			return false
		} else if isOk {
			return true
		}

		return a.contactDao.IsFriend(ctx, prem.UserId, prem.ReceiverId, false)
	} else if prem.TalkType == entity.ChatGroupMode {
		// 判断群是否解散
		group := &model.Group{}

		err := a.groupMemberDao.Db().First(group, "id = ?", prem.ReceiverId).Error
		if err != nil || group.Id == 0 || group.IsDismiss == 1 {
			return false
		}

		return a.groupMemberDao.IsMember(prem.ReceiverId, prem.UserId, true)
	}

	return false
}
