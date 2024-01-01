package logic

import (
	"context"
	"errors"
	"strings"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type MessageForwardLogic struct {
	db       *gorm.DB
	sequence *repo.Sequence
}

func NewMessageForwardLogic(db *gorm.DB, sequence *repo.Sequence) *MessageForwardLogic {
	return &MessageForwardLogic{db: db, sequence: sequence}
}

type ForwardRecord struct {
	MsgId      string
	ReceiverId int
	TalkType   int
}

// Verify 验证转发消息合法性
func (m *MessageForwardLogic) Verify(ctx context.Context, uid int, req *message.ForwardMessageRequest) error {

	query := m.db.WithContext(ctx).Model(&model.TalkRecords{})
	query.Where("msg_id in ?", req.MessageIds)

	if req.Receiver.TalkType == entity.ChatPrivateMode {
		subWhere := m.db.Where("user_id = ? and receiver_id = ?", uid, req.Receiver.ReceiverId)
		subWhere.Or("user_id = ? and receiver_id = ?", req.Receiver.ReceiverId, uid)
		query.Where(subWhere)
	}

	query.Where("talk_type = ?", req.Receiver.TalkType)
	query.Where("msg_type in ?", []int{1, 2, 3, 4, 5, 6, 7, 8, entity.ChatMsgTypeForward})
	query.Where("is_revoke = ?", 0)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if int(count) != len(req.MessageIds) {
		return errors.New("转发消息异常")
	}

	return nil
}

// MultiMergeForward 批量合并转发
func (m *MessageForwardLogic) MultiMergeForward(ctx context.Context, uid int, req *message.ForwardMessageRequest) ([]*ForwardRecord, error) {

	receives := make([]map[string]int, 0)

	for _, userId := range req.Uids {
		receives = append(receives, map[string]int{"receiver_id": int(userId), "talk_type": 1})
	}

	for _, gid := range req.Gids {
		receives = append(receives, map[string]int{"receiver_id": int(gid), "talk_type": 2})
	}

	tmpRecords, err := m.aggregation(ctx, req)
	if err != nil {
		return nil, err
	}

	extra := jsonutil.Encode(model.TalkRecordExtraForward{
		MsgIds:  req.MessageIds,
		Records: tmpRecords,
	})

	records := make([]*model.TalkRecords, 0, len(receives))
	for _, item := range receives {
		data := &model.TalkRecords{
			MsgId:      strutil.NewMsgId(),
			TalkType:   item["talk_type"],
			MsgType:    entity.ChatMsgTypeForward,
			UserId:     uid,
			ReceiverId: item["receiver_id"],
			Extra:      extra,
		}

		if data.TalkType == entity.ChatGroupMode {
			data.Sequence = m.sequence.Get(ctx, 0, data.ReceiverId)
		} else {
			data.Sequence = m.sequence.Get(ctx, uid, data.ReceiverId)
		}

		records = append(records, data)
	}

	if err := m.db.WithContext(ctx).Create(records).Error; err != nil {
		return nil, err
	}

	list := make([]*ForwardRecord, 0, len(records))
	for _, record := range records {
		list = append(list, &ForwardRecord{
			MsgId:      record.MsgId,
			ReceiverId: record.ReceiverId,
			TalkType:   record.TalkType,
		})
	}

	return list, nil
}

// MultiSplitForward 批量逐条转发
func (m *MessageForwardLogic) MultiSplitForward(ctx context.Context, uid int, req *message.ForwardMessageRequest) ([]*ForwardRecord, error) {
	var (
		receives = make([]map[string]int, 0)
		records  = make([]*model.TalkRecords, 0)
		db       = m.db.WithContext(ctx)
	)

	for _, userId := range req.Uids {
		receives = append(receives, map[string]int{"receiver_id": int(userId), "talk_type": model.TalkRecordTalkTypePrivate})
	}

	for _, gid := range req.Gids {
		receives = append(receives, map[string]int{"receiver_id": int(gid), "talk_type": model.TalkRecordTalkTypeGroup})
	}

	if err := db.Model(&model.TalkRecords{}).Where("id IN ?", req.MessageIds).Scan(&records).Error; err != nil {
		return nil, err
	}

	items := make([]*model.TalkRecords, 0, len(receives)*len(records))

	recordsLen := int64(len(records))
	for _, v := range receives {
		var sequences []int64

		if v["talk_type"] == model.TalkRecordTalkTypeGroup {
			sequences = m.sequence.BatchGet(ctx, 0, v["receiver_id"], recordsLen)
		} else {
			sequences = m.sequence.BatchGet(ctx, uid, v["receiver_id"], recordsLen)
		}

		for i, item := range records {
			items = append(items, &model.TalkRecords{
				MsgId:      strutil.NewMsgId(),
				TalkType:   v["talk_type"],
				MsgType:    item.MsgType,
				UserId:     uid,
				ReceiverId: v["receiver_id"],
				Sequence:   sequences[i],
				Extra:      item.Extra,
			})
		}
	}

	if err := db.Create(items).Error; err != nil {
		return nil, err
	}

	list := make([]*ForwardRecord, 0, len(items))
	for _, item := range items {
		list = append(list, &ForwardRecord{
			MsgId:      item.MsgId,
			ReceiverId: item.ReceiverId,
			TalkType:   item.TalkType,
		})
	}

	return list, nil
}

type forwardItem struct {
	MsgType  int    `json:"msg_type"`
	Extra    string `json:"extra"`
	Nickname string `json:"nickname"`
}

// 聚合转发数据
func (m *MessageForwardLogic) aggregation(ctx context.Context, req *message.ForwardMessageRequest) ([]map[string]any, error) {

	rows := make([]*forwardItem, 0, 3)

	query := m.db.WithContext(ctx).Table("talk_records").Select("talk_records.msg_type,talk_records.extra,users.nickname")
	query.Joins("left join users on users.id = talk_records.user_id")

	msgIds := req.MessageIds
	if len(msgIds) > 3 {
		msgIds = msgIds[:3]
	}

	query.Where("talk_records.msg_id in ?", msgIds)

	if err := query.Limit(3).Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0)
	for _, row := range rows {
		item := map[string]any{
			"nickname": row.Nickname,
		}

		switch row.MsgType {
		case entity.ChatMsgTypeText:
			extra := &model.TalkRecordExtraText{}
			if err := jsonutil.Decode(row.Extra, extra); err != nil {
				return nil, err
			}

			item["text"] = strutil.MtSubstr(strings.TrimSpace(extra.Content), 0, 30)
		case entity.ChatMsgTypeCode:
			item["text"] = "【代码消息】"
		case entity.ChatMsgTypeImage:
			item["text"] = "【图片消息】"
		case entity.ChatMsgTypeAudio:
			item["text"] = "【语音消息】"
		case entity.ChatMsgTypeVideo:
			item["text"] = "【视频消息】"
		case entity.ChatMsgTypeFile:
			item["text"] = "【文件消息】"
		case entity.ChatMsgTypeLocation:
			item["text"] = "【位置消息】"
		default:
			item["text"] = "【其它消息】"
		}

		data = append(data, item)
	}

	return data, nil
}
