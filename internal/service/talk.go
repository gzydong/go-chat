package service

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ ITalkService = (*TalkService)(nil)

type CollectOpt struct {
	UserId     int
	TalkType   int
	ReceiverId int
	MsgId      string
}

type ITalkService interface {
	Collect(ctx context.Context, opt *CollectOpt) error
	DeleteRecordList(ctx context.Context, opt *RemoveRecordListOpt) error
}

type TalkService struct {
	*repo.Source
	GroupMemberRepo *repo.GroupMember
}

type RemoveRecordListOpt struct {
	UserId     int
	TalkType   int
	ReceiverId int
	MsgIds     []string
}

// DeleteRecordList 删除消息记录
func (t *TalkService) DeleteRecordList(ctx context.Context, opt *RemoveRecordListOpt) error {
	var db = t.Source.Db().WithContext(ctx)

	// 私有消息直接删除
	if opt.TalkType == entity.ChatPrivateMode {
		err := db.Delete(model.TalkRecordFriend{}, "user_id = ? and msg_id in ?", opt.MsgIds).Error
		if err != nil {
			return err
		}
	}

	if !t.GroupMemberRepo.IsMember(ctx, opt.ReceiverId, opt.UserId, false) {
		return entity.ErrPermissionDenied
	}

	var findMsgIds []string
	db.Model(&model.TalkRecordGroup{}).
		Where("group_id = ? and msg_id in ?", opt.ReceiverId, opt.MsgIds, entity.ChatGroupMode).
		Pluck("msg_id", &findMsgIds)

	if len(opt.MsgIds) != len(findMsgIds) {
		return errors.New("删除异常! ")
	}

	items := make([]*model.TalkRecordGroupDel, 0, len(opt.MsgIds))
	for _, msgId := range opt.MsgIds {
		items = append(items, &model.TalkRecordGroupDel{
			MsgId:     msgId,
			UserId:    opt.UserId,
			CreatedAt: time.Now(),
		})
	}

	return db.Create(items).Error
}

// Collect 收藏表情包
func (t *TalkService) Collect(ctx context.Context, opt *CollectOpt) error {

	// 私有消息
	if opt.TalkType == entity.ChatPrivateMode {
		var record model.TalkRecordFriend

		err := t.Db().First(&record, "user_id = ? and friend_id = ? and msg_id = ?", opt.UserId, opt.ReceiverId, opt.MsgId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("消息ID不存在")
			}

			return err
		}

		if record.MsgType != entity.ChatMsgTypeImage {
			return errors.New("当前消息不支持收藏！")
		}

		var file model.TalkRecordExtraImage
		if err := jsonutil.Decode(record.Extra, &file); err != nil {
			return err
		}

		return t.Source.Db().Create(&model.EmoticonItem{
			UserId: opt.UserId,
			Url:    file.Url,
		}).Error
	}

	if !t.GroupMemberRepo.IsMember(ctx, opt.ReceiverId, opt.UserId, true) {
		return entity.ErrPermissionDenied
	}

	var record model.TalkRecordGroup
	if err := t.Source.Db().First(&record, "group_id = ? and msg_id = ?", opt.MsgId).Error; err != nil {
		return err
	}

	if record.MsgType != entity.ChatMsgTypeImage {
		return errors.New("当前消息不支持收藏！")
	}

	var file model.TalkRecordExtraImage
	if err := jsonutil.Decode(record.Extra, &file); err != nil {
		return err
	}

	return t.Source.Db().Create(&model.EmoticonItem{
		UserId: opt.UserId,
		Url:    file.Url,
	}).Error
}
