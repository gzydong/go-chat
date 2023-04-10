package service

import (
	"context"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
)

type AuthPermissionService struct {
	contactRepo     *repo.Contact
	groupMemberRepo *repo.GroupMember
	organizeRepo    *organize.Organize
}

func NewAuthPermissionService(contactRepo *repo.Contact, groupMemberRepo *repo.GroupMember, organizeRepo *organize.Organize) *AuthPermissionService {
	return &AuthPermissionService{contactRepo: contactRepo, groupMemberRepo: groupMemberRepo, organizeRepo: organizeRepo}
}

type AuthPermission struct {
	TalkType   int
	UserId     int
	ReceiverId int
}

func (a *AuthPermissionService) IsAuth(ctx context.Context, params *AuthPermission) bool {
	if params.TalkType == entity.ChatPrivateMode {
		// 这里需要判断双方是否都是企业成员，如果是则无需添加好友即可聊天
		if isOk, err := a.organizeRepo.IsQiyeMember(ctx, params.UserId, params.ReceiverId); err != nil {
			logger.Error("[AuthPermission IsAuth] 查询数据异常 err: ", err)
			return false
		} else if isOk {
			return true
		}

		return a.contactRepo.IsFriend(ctx, params.UserId, params.ReceiverId, false)
	} else if params.TalkType == entity.ChatGroupMode {
		var group model.Group
		err := a.groupMemberRepo.Db.WithContext(ctx).First(&group, "id = ?", params.ReceiverId).Error
		if err != nil || group.Id == 0 || group.IsDismiss == 1 {
			return false
		}

		return a.groupMemberRepo.IsMember(ctx, params.ReceiverId, params.UserId, true)
	}

	return false
}
