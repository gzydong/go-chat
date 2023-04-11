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
	organize *organize.Organize
	contact  *repo.Contact
}

func NewAuthService(organize *organize.Organize, contact *repo.Contact) *AuthService {
	return &AuthService{organize: organize, contact: contact}
}

type AuthOption struct {
	TalkType   int
	UserId     int
	ReceiverId int
}

func (t *AuthService) IsAuth(ctx context.Context, opt *AuthOption) error {

	if opt.TalkType == entity.ChatPrivateMode {
		// 这里需要判断双方是否都是企业成员，如果是则无需添加好友即可聊天
		if isOk, err := t.organize.IsQiyeMember(ctx, opt.UserId, opt.ReceiverId); err != nil {
			return errors.New("系统繁忙，请稍后再试！！！")
		} else if isOk {
			return nil
		}

		if t.contact.IsFriend(ctx, opt.UserId, opt.ReceiverId, false) {
			return nil
		}

		return errors.New("暂无权限发送消息！")
	}

	var memberInfo model.GroupMember
	err := t.contact.Db.First(memberInfo, "group_id = ? and user_id = ?", opt.ReceiverId, opt.UserId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("暂无权限发送消息！")
		}

		return errors.New("系统繁忙，请稍后再试！！！")
	}

	if memberInfo.IsQuit == 1 {
		return errors.New("暂无权限发送消息！")
	}

	if memberInfo.IsMute == 1 {
		return errors.New("已被群主或管理员禁言！")
	}

	return nil
}
