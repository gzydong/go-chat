package service

import (
	"context"

	"go-chat/internal/dao"
	"go-chat/internal/entity"
)

type AuthPermissionService struct {
	contactDao     *dao.ContactDao
	groupMemberDao *dao.GroupMemberDao
}

func NewAuthPermissionService(contactDao *dao.ContactDao, groupMemberDao *dao.GroupMemberDao) *AuthPermissionService {
	return &AuthPermissionService{contactDao: contactDao, groupMemberDao: groupMemberDao}
}

type AuthPermission struct {
	TalkType   int
	UserId     int
	ReceiverId int
}

func (a *AuthPermissionService) IsAuth(ctx context.Context, prem *AuthPermission) bool {
	if prem.TalkType == entity.ChatPrivateMode {
		return a.contactDao.IsFriend(ctx, prem.UserId, prem.ReceiverId, true)
	} else if prem.TalkType == entity.ChatGroupMode {
		return a.groupMemberDao.IsMember(prem.ReceiverId, prem.UserId, true)
	}

	return false
}
