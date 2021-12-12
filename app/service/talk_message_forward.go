package service

import (
	"context"
	"go-chat/app/entity"
	"go-chat/app/model"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"gorm.io/gorm"
	"strings"
)

type ForwardParams struct {
	UserId     int   `json:"user_id"`
	ReceiverId int   `json:"receiver_id"`
	TalkType   int   `json:"talk_type"`
	RecordsIds []int `json:"records_ids"`
	UserIds    []int `json:"user_ids"`
	GroupIds   []int `json:"group_ids"`
}

type TalkMessageForwardService struct {
	*BaseService
}

func NewTalkMessageForwardService(base *BaseService) *TalkMessageForwardService {
	return &TalkMessageForwardService{base}
}

// 验证消息转发
func (t *TalkMessageForwardService) verifyForward(forward *ForwardParams) error {
	return nil
}

// SendForwardMessage 推送转发消息
func (t *TalkMessageForwardService) SendForwardMessage(ctx context.Context, forward *ForwardParams) error {
	var (
		err   error
		items []*PushReceive
	)

	if err = t.verifyForward(forward); err != nil {
		return err
	}

	if forward.TalkType == 1 {
		items, err = t.MultiMergeForward(ctx, forward)
	} else {
		items, err = t.MultiSplitForward(ctx, forward)
	}

	for _, item := range items {
		body := entity.JsonText{
			"event": entity.EventTalk,
			"data": entity.JsonText{
				"sender_id":   int64(forward.UserId),
				"receiver_id": int64(item.ReceiverId),
				"talk_type":   item.TalkType,
				"record_id":   int64(item.RecordId),
			}.Json(),
		}

		t.rds.Publish(ctx, entity.SubscribeWsGatewayAll, body.Json())
	}

	return nil
}

type Receives struct {
	ReceiverId int `json:"receiver_id"`
	TalkType   int `json:"talk_type"`
}

type PushReceive struct {
	RecordId   int `json:"record_id"`
	ReceiverId int `json:"receiver_id"`
	TalkType   int `json:"talk_type"`
}

type ForwardMsgItem struct {
	MsgType  int    `json:"msg_type"`
	Content  string `json:"content"`
	Nickname string `json:"nickname"`
}

// 聚合转发数据
func (t *TalkMessageForwardService) aggregation(ctx context.Context, forward *ForwardParams) (string, error) {
	rows := make([]*ForwardMsgItem, 0)
	query := t.db.Table("talk_records")
	query.Joins("left join users on users.id = talk_records.user_id")
	query.Where("talk_records.id in ?", forward.RecordsIds[:3])

	if err := query.Limit(3).Scan(&rows).Error; err != nil {
		return "", err
	}

	data := make([]map[string]interface{}, 0)
	for _, row := range rows {
		item := map[string]interface{}{}

		switch row.MsgType {
		case entity.MsgTypeText:
			text := strings.TrimSpace(row.Content)
			item["nickname"] = row.Nickname
			item["text"] = strutil.MtSubstr(&text, 0, 30)
		case entity.MsgTypeCode:
			item["nickname"] = row.Nickname
			item["text"] = "【代码消息】"
		case entity.MsgTypeFile:
			item["nickname"] = row.Nickname
			item["text"] = "【文件消息】"
		}

		data = append(data, item)
	}

	return jsonutil.JsonEncode(data), nil
}

// MultiMergeForward 转发消息（多条合并转发）
func (t *TalkMessageForwardService) MultiMergeForward(ctx context.Context, forward *ForwardParams) ([]*PushReceive, error) {
	var (
		receives = make([]*Receives, 0)
		arr      = make([]*PushReceive, 0)
	)

	for _, uid := range forward.UserIds {
		receives = append(receives, &Receives{uid, 1})
	}

	for _, gid := range forward.GroupIds {
		receives = append(receives, &Receives{gid, 2})
	}

	text, err := t.aggregation(ctx, forward)
	if err != nil {
		return nil, err
	}

	str := slice.IntToIds(forward.RecordsIds)
	err = t.db.Transaction(func(tx *gorm.DB) error {
		forwards := make([]*model.TalkRecordsForward, 0)
		records := make([]*model.TalkRecords, 0)

		for _, receive := range receives {
			records = append(records, &model.TalkRecords{
				TalkType:   receive.TalkType,
				MsgType:    entity.MsgTypeForward,
				UserId:     forward.UserId,
				ReceiverId: receive.ReceiverId,
			})
		}

		if err := tx.Create(records).Error; err != nil {
			return err
		}

		for _, record := range records {
			forwards = append(forwards, &model.TalkRecordsForward{
				RecordId:  record.Id,
				UserId:    record.UserId,
				RecordsId: str,
				Text:      text,
			})

			arr = append(arr, &PushReceive{
				RecordId:   record.Id,
				ReceiverId: record.ReceiverId,
				TalkType:   record.TalkType,
			})
		}

		if err := tx.Create(&forwards).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return arr, nil
}

// MultiSplitForward 转发消息（多条拆分转发）
func (t *TalkMessageForwardService) MultiSplitForward(ctx context.Context, forward *ForwardParams) ([]*PushReceive, error) {
	var (
		receives = make([]*Receives, 0)
		arr      = make([]*PushReceive, 0)
		records  = make([]*model.TalkRecords, 0)
	)

	for _, uid := range forward.UserIds {
		receives = append(receives, &Receives{uid, 1})
	}

	for _, gid := range forward.GroupIds {
		receives = append(receives, &Receives{gid, 2})
	}

	if err := t.db.Model(&model.TalkRecords{}).Where("id = ?", forward.RecordsIds).Scan(&records).Error; err != nil {
		return nil, err
	}

	err := t.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range records {
			items := make([]*model.TalkRecords, 0)
			for _, receive := range receives {
				items = append(items, &model.TalkRecords{
					TalkType:   receive.TalkType,
					MsgType:    item.MsgType,
					UserId:     forward.UserId,
					ReceiverId: receive.ReceiverId,
					Content:    item.Content,
				})
			}

			if err := tx.Create(items).Error; err != nil {
				return err
			}

			files := make([]model.TalkRecordsFile, 0)
			codes := make([]model.TalkRecordsCode, 0)

			for _, record := range items {
				arr = append(arr, &PushReceive{
					RecordId:   record.Id,
					ReceiverId: record.ReceiverId,
					TalkType:   record.TalkType,
				})
			}

			if len(files) > 0 {

			}

			if len(codes) > 0 {

			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return arr, nil
}
