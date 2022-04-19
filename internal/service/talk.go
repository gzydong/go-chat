package service

import (
	"context"
	"errors"
	"time"

	"go-chat/internal/dao"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/sliceutil"
)

type TalkMessageDeleteOpts struct {
	UserId     int
	TalkType   int
	ReceiverId int
	RecordIds  string
}

type TalkService struct {
	*BaseService
	groupMemberDao *dao.GroupMemberDao
}

func NewTalkService(baseService *BaseService, groupMemberDao *dao.GroupMemberDao) *TalkService {
	return &TalkService{BaseService: baseService, groupMemberDao: groupMemberDao}
}

// RemoveRecords 删除消息记录
func (s *TalkService) RemoveRecords(ctx context.Context, opts *TalkMessageDeleteOpts) error {

	// 需要删除的消息记录ID
	ids := sliceutil.UniqueInt(sliceutil.ParseIds(opts.RecordIds))

	// 查询的ids
	findIds := make([]int64, 0)

	if opts.TalkType == entity.ChatPrivateMode {
		subQuery := s.db.Where("user_id = ? and receiver_id = ?", opts.UserId, opts.ReceiverId).Or("user_id = ? and receiver_id = ?", opts.ReceiverId, opts.UserId)

		s.db.Model(&model.TalkRecords{}).Where("id in ?", ids).Where("talk_type = ?", entity.ChatPrivateMode).Where(subQuery).Pluck("id", &findIds)
	} else {
		if !s.groupMemberDao.IsMember(opts.ReceiverId, opts.UserId, false) {
			return entity.ErrPermissionDenied
		}

		s.db.Model(&model.TalkRecords{}).Where("id in ? and talk_type = ?", ids, entity.ChatGroupMode).Pluck("id", &findIds)
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

	return s.db.Create(items).Error
}

// CollectRecord 收藏表情包
func (s *TalkService) CollectRecord(ctx context.Context, uid int, recordId int) error {
	var (
		err      error
		record   model.TalkRecords
		fileInfo model.TalkRecordsFile
	)

	if err = s.db.First(&record, recordId).Error; err != nil {
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
		if !s.groupMemberDao.IsMember(record.ReceiverId, uid, true) {
			return entity.ErrPermissionDenied
		}
	}

	if err = s.db.First(&fileInfo, "record_id = ? and type = ?", record.Id, 1).Error; err != nil {
		return err
	}

	emoticon := &model.EmoticonItem{
		UserId:     uid,
		Url:        fileInfo.Url,
		FileSuffix: fileInfo.Suffix,
		FileSize:   fileInfo.Size,
	}

	return s.db.Create(emoticon).Error
}
