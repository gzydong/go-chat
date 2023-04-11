package service

import (
	"context"
	"errors"

	"go-chat/internal/entity"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
	"gorm.io/gorm"
)

type AuthService struct {
	organize    *organize.Organize
	contact     *repo.Contact
	groupMember *repo.GroupMember
}

func NewAuthService(organize *organize.Organize, contact *repo.Contact, groupMember *repo.GroupMember) *AuthService {
	return &AuthService{organize: organize, contact: contact, groupMember: groupMember}
}

type AuthOption struct {
	TalkType   int
	UserId     int
	ReceiverId int
}

func (a *AuthService) IsAuth(ctx context.Context, opt *AuthOption) error {

	if opt.TalkType == entity.ChatPrivateMode {
		if isOk, err := a.organize.IsQiyeMember(ctx, opt.UserId, opt.ReceiverId); err != nil {
			return errors.New("系统繁忙，请稍后再试！！！")
		} else if isOk {
			return nil
		}

		if a.contact.IsFriend(ctx, opt.UserId, opt.ReceiverId, false) {
			return nil
		}

		return errors.New("暂无权限发送消息！")
	}

	memberInfo, err := a.groupMember.FindByUserId(ctx, opt.ReceiverId, opt.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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

	return nil
}
