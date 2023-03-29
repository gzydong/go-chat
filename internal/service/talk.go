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

type TalkService struct {
	*repo.Source
	groupMemberRepo *repo.GroupMember
}

func NewTalkService(source *repo.Source, groupMemberRepo *repo.GroupMember) *TalkService {
	return &TalkService{Source: source, groupMemberRepo: groupMemberRepo}
}

type TalkMessageDeleteOpt struct {
	UserId     int
	TalkType   int
	ReceiverId int
	RecordIds  string
}

// RemoveRecords 删除消息记录
func (s *TalkService) RemoveRecords(ctx context.Context, opts *TalkMessageDeleteOpt) error {

	// 需要删除的消息记录ID
	ids := sliceutil.Unique(sliceutil.ParseIds(opts.RecordIds))

	// 查询的ids
	findIds := make([]int64, 0)

	if opts.TalkType == entity.ChatPrivateMode {
		subQuery := s.Db().Where("user_id = ? and receiver_id = ?", opts.UserId, opts.ReceiverId).Or("user_id = ? and receiver_id = ?", opts.ReceiverId, opts.UserId)

		s.Db().Model(&model.TalkRecords{}).Where("id in ?", ids).Where("talk_type = ?", entity.ChatPrivateMode).Where(subQuery).Pluck("id", &findIds)
	} else {
		if !s.groupMemberRepo.IsMember(ctx, opts.ReceiverId, opts.UserId, false) {
			return entity.ErrPermissionDenied
		}

		s.Db().Model(&model.TalkRecords{}).Where("id in ? and talk_type = ?", ids, entity.ChatGroupMode).Pluck("id", &findIds)
	}

	if len(ids) != len(findIds) {
		return errors.New("删除异常! ")
	}

	items := make([]*model.TalkRecordsDelete, 0)
	for _, val := range ids {
		items = append(items, &model.TalkRecordsDelete{
			RecordId:  val,
			UserId:    opts.UserId,
			CreatedAt: time.Now(),
		})
	}

	return s.Db().Create(items).Error
}

// CollectRecord 收藏表情包
func (s *TalkService) CollectRecord(ctx context.Context, uid int, recordId int) error {

	var record model.TalkRecords
	if err := s.Db().First(&record, recordId).Error; err != nil {
		return err
	}

	if record.MsgType != entity.MsgTypeFile {
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
		if !s.groupMemberRepo.IsMember(ctx, record.ReceiverId, uid, true) {
			return entity.ErrPermissionDenied
		}
	}

	var fileInfo model.TalkRecordExtraFile
	if err := jsonutil.Decode(record.Extra, &fileInfo); err != nil {
		return err
	}

	return s.Db().Create(&model.EmoticonItem{
		UserId:     uid,
		Url:        fileInfo.Url,
		FileSuffix: fileInfo.Suffix,
		FileSize:   fileInfo.Size,
	}).Error
}
