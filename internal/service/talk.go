package service

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

var _ ITalkService = (*TalkService)(nil)

type ITalkService interface {
	Collect(ctx context.Context, uid int, recordId int) error
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
	RecordIds  string
}

// DeleteRecordList 删除消息记录
func (t *TalkService) DeleteRecordList(ctx context.Context, opt *RemoveRecordListOpt) error {

	var (
		db      = t.Source.Db().WithContext(ctx)
		findIds []int64
		ids     = sliceutil.Unique(sliceutil.ParseIds(opt.RecordIds))
	)

	if opt.TalkType == entity.ChatPrivateMode {
		subQuery := db.Where("user_id = ? and receiver_id = ?", opt.UserId, opt.ReceiverId).Or("user_id = ? and receiver_id = ?", opt.ReceiverId, opt.UserId)
		db.Model(&model.TalkRecords{}).Where("id in ?", ids).Where("talk_type = ?", entity.ChatPrivateMode).Where(subQuery).Pluck("id", &findIds)
	} else {
		if !t.GroupMemberRepo.IsMember(ctx, opt.ReceiverId, opt.UserId, false) {
			return entity.ErrPermissionDenied
		}

		db.Model(&model.TalkRecords{}).Where("id in ? and talk_type = ?", ids, entity.ChatGroupMode).Pluck("id", &findIds)
	}

	if len(ids) != len(findIds) {
		return errors.New("删除异常! ")
	}

	items := make([]*model.TalkRecordsDelete, 0, len(ids))
	for _, val := range ids {
		items = append(items, &model.TalkRecordsDelete{
			RecordId:  val,
			UserId:    opt.UserId,
			CreatedAt: time.Now(),
		})
	}

	return db.Create(items).Error
}

// Collect 收藏表情包
func (t *TalkService) Collect(ctx context.Context, uid int, recordId int) error {

	var record model.TalkRecords
	if err := t.Source.Db().First(&record, recordId).Error; err != nil {
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
		UserId:   uid,
		Url:      file.Url,
		FileSize: file.Size,
	}).Error
}
