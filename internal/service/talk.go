package service

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

var _ ITalkService = (*TalkService)(nil)

type ITalkService interface {
	Collect(ctx context.Context, uid int, msgId string) error
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

	var (
		db         = t.Source.Db().WithContext(ctx)
		findMsgIds []string
	)

	if opt.TalkType == entity.ChatPrivateMode {
		subQuery := db.Where("user_id = ? and receiver_id = ?", opt.UserId, opt.ReceiverId).Or("user_id = ? and receiver_id = ?", opt.ReceiverId, opt.UserId)
		db.Model(&model.TalkRecords{}).Where("msg_id in ?", opt.MsgIds).Where("talk_type = ?", entity.ChatPrivateMode).Where(subQuery).Pluck("msg_id", &findMsgIds)
	} else {
		if !t.GroupMemberRepo.IsMember(ctx, opt.ReceiverId, opt.UserId, false) {
			return entity.ErrPermissionDenied
		}

		db.Model(&model.TalkRecords{}).Where("msg_id in ? and talk_type = ?", opt.MsgIds, entity.ChatGroupMode).Pluck("msg_id", &findMsgIds)
	}

	if len(opt.MsgIds) != len(findMsgIds) {
		return errors.New("删除异常! ")
	}

	items := make([]*model.TalkRecordsDelete, 0, len(opt.MsgIds))
	for _, msgId := range opt.MsgIds {
		items = append(items, &model.TalkRecordsDelete{
			MsgId:     msgId,
			UserId:    opt.UserId,
			CreatedAt: time.Now(),
		})
	}

	return db.Create(items).Error
}

// Collect 收藏表情包
func (t *TalkService) Collect(ctx context.Context, uid int, msgId string) error {

	var record model.TalkRecords
	if err := t.Source.Db().First(&record, "msg_id = ?", msgId).Error; err != nil {
		return err
	}

	if record.MsgType != entity.ChatMsgTypeImage {
		return errors.New("当前消息不支持收藏！")
	}

	if record.IsRevoke == 1 {
		return errors.New("当前消息不支持收藏！")
	}

	if record.TalkType == entity.ChatPrivateMode {
		if record.UserId != uid && record.ReceiverId != uid {
			return entity.ErrPermissionDenied
		}
	} else if record.TalkType == entity.ChatGroupMode {
		if !t.GroupMemberRepo.IsMember(ctx, record.ReceiverId, uid, true) {
			return entity.ErrPermissionDenied
		}
	}

	var file model.TalkRecordExtraImage
	if err := jsonutil.Decode(record.Extra, &file); err != nil {
		return err
	}

	return t.Source.Db().Create(&model.EmoticonItem{
		UserId:    uid,
		Url:       file.Url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}
