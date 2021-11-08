package service

import (
	"context"
	"errors"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/slice"
	"time"
)

type TalkService struct {
	*BaseService
	groupMemberService *GroupMemberService
}

func NewTalkService(base *BaseService, groupMemberService *GroupMemberService) *TalkService {
	return &TalkService{base, groupMemberService}
}

// RemoveRecords 删除消息记录
// @params uid 用户ID
// @params req 请求参数
func (s *TalkService) RemoveRecords(ctx context.Context, uid int, req *request.DeleteMessageRequest) error {

	// 需要删除的消息记录ID
	ids := slice.UniqueInt(slice.ParseIds(req.RecordIds))

	// 查询的ids
	findIds := make([]int64, 0)

	if req.TalkType == entity.PrivateChat {
		s.db.Model(model.TalkRecords{}).Where("id in ?", ids).Where("talk_type = ?", entity.PrivateChat).Where(
			s.db.Where("user_id = ? and receiver_id = ?", uid, req.ReceiverId).
				Or("user_id = ? and receiver_id = ?", req.ReceiverId, uid)).Pluck("id", &findIds)
	} else {
		if !s.groupMemberService.IsMember(req.ReceiverId, uid) {
			return errors.New("非群成员，暂无权限! ")
		}

		s.db.Model(model.TalkRecords{}).Where("id in ? and talk_type = ?", ids, entity.GroupChat).Pluck("id", &findIds)
	}

	if len(ids) != len(findIds) {
		return errors.New("删除异常! ")
	}

	var items []*model.TalkRecordsDelete
	for _, val := range ids {
		items = append(items, &model.TalkRecordsDelete{
			RecordId:  val,
			UserId:    uid,
			CreatedAt: time.Now(),
		})
	}

	return s.db.Create(items).Error
}

// CollectRecord 收藏表情包
// @params uid      用户ID
// @params recordId 消息记录ID
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

	if record.TalkType == entity.PrivateChat {
		if record.UserId != uid && record.ReceiverId != uid {
			return errors.New("暂无权限收藏！")
		}
	} else if record.TalkType == entity.GroupChat {
		if !s.groupMemberService.IsMember(record.ReceiverId, uid) {
			return errors.New("暂无权限收藏！")
		}
	}

	if err = s.db.First(&fileInfo, "record_id = ? and file_type = ?", record.ID, 1).Error; err != nil {
		return err
	}

	emoticon := &model.EmoticonItem{
		EmoticonId: 0,
		UserId:     uid,
		Url:        fileInfo.SaveDir,
		FileSuffix: fileInfo.FileSuffix,
		FileSize:   fileInfo.FileSize,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.db.Create(emoticon).Error
}
