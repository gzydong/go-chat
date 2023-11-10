package service

import (
	"context"
	"errors"

	"go-chat/internal/entity"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IAuthService = (*AuthService)(nil)

type IAuthService interface {
	IsAuth(ctx context.Context, opt *AuthOption) error
}

type AuthService struct {
	OrganizeRepo    *repo.Organize
	ContactRepo     *repo.Contact
	GroupRepo       *repo.Group
	GroupMemberRepo *repo.GroupMember
}

type AuthOption struct {
	TalkType          int
	UserId            int
	ReceiverId        int
	IsVerifyGroupMute bool
}

func (a *AuthService) IsAuth(ctx context.Context, opt *AuthOption) error {

	if opt.TalkType == entity.ChatPrivateMode {
		if isOk, err := a.OrganizeRepo.IsQiyeMember(ctx, opt.UserId, opt.ReceiverId); err != nil {
			return errors.New("系统繁忙，请稍后再试！！！")
		} else if isOk {
			return nil
		}

		if a.ContactRepo.IsFriend(ctx, opt.UserId, opt.ReceiverId, false) {
			return nil
		}

		return errors.New("暂无权限发送消息！")
	}

	groupInfo, err := a.GroupRepo.FindById(ctx, opt.ReceiverId)
	if err != nil {
		return err
	}

	if groupInfo.IsDismiss == 1 {
		return errors.New("此群聊已解散！")
	}

	memberInfo, err := a.GroupMemberRepo.FindByUserId(ctx, opt.ReceiverId, opt.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("暂无权限发送消息！")
		}

		return errors.New("系统繁忙，请稍后再试！！！")
	}

	if memberInfo.IsQuit == model.GroupMemberQuitStatusYes {
		return errors.New("暂无权限发送消息！")
	}

	if memberInfo.IsMute == model.GroupMemberMuteStatusYes {
		return errors.New("已被群主或管理员禁言！")
	}

	if opt.IsVerifyGroupMute && groupInfo.IsMute == 1 && memberInfo.Leader == 0 {
		return errors.New("此群聊已开启全员禁言！")
	}

	return nil
}
