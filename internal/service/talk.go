package service

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ ITalkService = (*TalkService)(nil)

type TalkRevokeOption struct {
	UserId   int
	TalkMode int
	MsgId    string
}

type TalkDeleteRecordOption struct {
	UserId   int
	TalkMode int
	ToFromId int
	MsgIds   []string
}

type ITalkService interface {
	DeleteRecord(ctx context.Context, opt *TalkDeleteRecordOption) error
	Revoke(ctx context.Context, opt *TalkRevokeOption) error
}

type TalkService struct {
	*repo.Source
	GroupMemberRepo *repo.GroupMember
}

// DeleteRecord 删除消息记录
func (t *TalkService) DeleteRecord(ctx context.Context, opt *TalkDeleteRecordOption) error {
	var db = t.Source.Db().WithContext(ctx)

	// 私有消息直接更新删除状态
	if opt.TalkMode == entity.ChatPrivateMode {
		return db.Model(model.TalkUserMessage{}).
			Where("user_id = ? and msg_id in ?", opt.UserId, opt.MsgIds).
			Update("is_deleted", model.Yes).Error
	}

	if !t.GroupMemberRepo.IsMember(ctx, opt.ToFromId, opt.UserId, false) {
		return entity.ErrPermissionDenied
	}

	var findMsgIds []string
	db.Model(&model.TalkGroupMessage{}).
		Where("group_id = ? and msg_id in ?", opt.ToFromId, opt.MsgIds).
		Pluck("msg_id", &findMsgIds)

	if len(opt.MsgIds) != len(findMsgIds) {
		return errors.New("删除异常! ")
	}

	items := make([]*model.TalkGroupMessageDel, 0, len(opt.MsgIds))
	for _, msgId := range opt.MsgIds {
		items = append(items, &model.TalkGroupMessageDel{
			MsgId:     msgId,
			GroupId:   opt.ToFromId,
			UserId:    opt.UserId,
			CreatedAt: time.Now(),
		})
	}

	// 删除后清除最后一条记录
	return db.Create(items).Error
}

// Revoke 撤回消息
// TODO 撤回消息后需要推送消息
func (t *TalkService) Revoke(ctx context.Context, opt *TalkRevokeOption) error {

	db := t.Db().WithContext(ctx)

	switch opt.TalkMode {
	case entity.ChatPrivateMode:
		var record model.TalkUserMessage

		err := db.First(&record, "msg_id = ? and from_id = ?", opt.MsgId, opt.UserId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("消息ID不存在")
			}

			return err
		}

		if record.IsRevoked == model.Yes {
			return errors.New("消息已撤回")
		}

		if time.Now().Unix() > record.SendTime.Add(3*time.Minute).Unix() {
			return errors.New("超出有效撤回时间范围，无法进行撤销！")
		}

		return db.Model(&model.TalkUserMessage{}).
			Where("org_msg_id = ?", record.OrgMsgId).
			Update("is_revoked", model.Yes).Error

	case entity.ChatGroupMode:
		var record model.TalkGroupMessage

		err := db.First(&record, "msg_id = ? and from_id = ?", opt.MsgId, opt.UserId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("消息ID不存在")
			}

			return err
		}

		if record.IsRevoked == model.Yes {
			return errors.New("消息已撤回")
		}

		if time.Now().Unix() > record.SendTime.Add(3*time.Minute).Unix() {
			return errors.New("超出有效撤回时间范围，无法进行撤销！")
		}

		return db.Model(&model.TalkGroupMessage{}).
			Where("msg_id = ?", record.MsgId).
			Update("is_revoked", model.Yes).Error
	}

	return errors.New("暂不支持撤回消息")
}
