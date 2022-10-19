package logic

import (
	"context"
	"errors"
	"strings"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type MessageForwardLogic struct {
	db       *gorm.DB
	sequence *cache.Sequence
}

func NewMessageForwardLogic(db *gorm.DB, sequence *cache.Sequence) *MessageForwardLogic {
	return &MessageForwardLogic{db: db, sequence: sequence}
}

type ForwardRecord struct {
	RecordId   int
	ReceiverId int
	TalkType   int
}

// Verify 验证转发消息合法性
func (m *MessageForwardLogic) Verify(ctx context.Context, uid int, req *message.ForwardMessageRequest) error {

	query := m.db.Model(&model.TalkRecords{})
	query.Where("id in ?", req.MessageIds)

	if req.Receiver.TalkType == entity.ChatPrivateMode {
		subWhere := m.db.Where("user_id = ? and receiver_id = ?", uid, req.Receiver.ReceiverId)
		subWhere.Or("user_id = ? and receiver_id = ?", req.Receiver.ReceiverId, uid)
		query.Where(subWhere)
	}

	query.Where("talk_type = ?", req.Receiver.TalkType)
	query.Where("msg_type in ?", []int{1, 2, 4})
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
	var (
		receives = make([]map[string]int, 0)
		arr      = make([]*ForwardRecord, 0)
	)

	for _, uid := range req.Uids {
		receives = append(receives, map[string]int{
			"id":   int(uid),
			"type": 1,
		})
	}

	for _, gid := range req.Gids {
		receives = append(receives, map[string]int{
			"id":   int(gid),
			"type": 2,
		})
	}

	text, err := m.aggregation(ctx, req)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0)
	for _, id := range req.MessageIds {
		ids = append(ids, int(id))
	}

	str := sliceutil.ToIds(ids)
	err = m.db.Transaction(func(tx *gorm.DB) error {
		forwards := make([]*model.TalkRecordsForward, 0, len(receives))
		records := make([]*model.TalkRecords, 0, len(receives))

		for _, item := range receives {

			data := &model.TalkRecords{
				MsgId:      strutil.NewUuid(),
				TalkType:   item["type"],
				MsgType:    entity.MsgTypeForward,
				UserId:     uid,
				ReceiverId: item["id"],
			}

			if data.TalkType == entity.ChatGroupMode {
				data.Sequence = m.sequence.Seq(ctx, 0, data.ReceiverId)
			} else {
				data.Sequence = m.sequence.Seq(ctx, uid, data.ReceiverId)
			}

			records = append(records, data)
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

			arr = append(arr, &ForwardRecord{
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

// MultiSplitForward 批量逐条转发
func (m *MessageForwardLogic) MultiSplitForward(ctx context.Context, uid int, req *message.ForwardMessageRequest) ([]*ForwardRecord, error) {
	var (
		receives  = make([]map[string]int, 0)
		arr       = make([]*ForwardRecord, 0)
		records   = make([]*model.TalkRecords, 0)
		hashFiles = make(map[int]*model.TalkRecordsFile)
		hashCodes = make(map[int]*model.TalkRecordsCode)
	)

	for _, uid := range req.Uids {
		receives = append(receives, map[string]int{
			"id":   int(uid),
			"type": 1,
		})
	}

	for _, gid := range req.Gids {
		receives = append(receives, map[string]int{
			"id":   int(gid),
			"type": 2,
		})
	}

	if err := m.db.Model(&model.TalkRecords{}).Where("id IN ?", req.MessageIds).Scan(&records).Error; err != nil {
		return nil, err
	}

	codeIds, fileIds := make([]int, 0), make([]int, 0)

	for _, record := range records {
		switch record.MsgType {
		case entity.MsgTypeFile:
			fileIds = append(fileIds, record.Id)
		case entity.MsgTypeCode:
			codeIds = append(codeIds, record.Id)
		}
	}

	if len(codeIds) > 0 {
		items := make([]*model.TalkRecordsCode, 0)
		if err := m.db.Model(&model.TalkRecordsCode{}).Where("record_id in ?", codeIds).Scan(&items).Error; err == nil {
			for i := range items {
				hashCodes[items[i].RecordId] = items[i]
			}
		}
	}

	if len(fileIds) > 0 {
		items := make([]*model.TalkRecordsFile, 0)
		if err := m.db.Model(&model.TalkRecordsFile{}).Where("record_id in ?", fileIds).Scan(&items).Error; err == nil {
			for i := range items {
				hashFiles[items[i].RecordId] = items[i]
			}
		}
	}

	err := m.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range records {
			items := make([]*model.TalkRecords, 0, len(receives))
			files := make([]*model.TalkRecordsFile, 0)
			codes := make([]*model.TalkRecordsCode, 0)

			for _, v := range receives {
				data := &model.TalkRecords{
					MsgId:      strutil.NewUuid(),
					TalkType:   v["type"],
					MsgType:    item.MsgType,
					UserId:     uid,
					ReceiverId: v["id"],
					Content:    item.Content,
				}

				if data.TalkType == entity.ChatGroupMode {
					data.Sequence = m.sequence.Seq(ctx, 0, data.ReceiverId)
				} else {
					data.Sequence = m.sequence.Seq(ctx, uid, data.ReceiverId)
				}

				items = append(items, data)
			}

			if err := tx.Create(items).Error; err != nil {
				return err
			}

			for _, record := range items {
				arr = append(arr, &ForwardRecord{
					RecordId:   record.Id,
					ReceiverId: record.ReceiverId,
					TalkType:   record.TalkType,
				})

				switch record.MsgType {
				case entity.MsgTypeFile:
					if file, ok := hashFiles[item.Id]; ok {
						files = append(files, &model.TalkRecordsFile{
							RecordId:     record.Id,
							UserId:       uid,
							Source:       file.Source,
							Type:         file.Type,
							Drive:        file.Drive,
							OriginalName: file.OriginalName,
							Suffix:       file.Suffix,
							Size:         file.Size,
							Path:         file.Path,
						})
					}
				case entity.MsgTypeCode:
					if code, ok := hashCodes[item.Id]; ok {
						codes = append(codes, &model.TalkRecordsCode{
							RecordId: record.Id,
							UserId:   uid,
							Lang:     code.Lang,
							Code:     code.Code,
						})
					}
				}
			}

			if len(files) > 0 {
				if err := tx.Create(files).Error; err != nil {
					return err
				}
			}

			if len(codes) > 0 {
				if err := tx.Create(codes).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return arr, nil
}

type forwardItem struct {
	MsgType  int    `json:"msg_type"`
	Content  string `json:"content"`
	Nickname string `json:"nickname"`
}

// 聚合转发数据
func (m *MessageForwardLogic) aggregation(ctx context.Context, req *message.ForwardMessageRequest) (string, error) {

	rows := make([]*forwardItem, 0)

	query := m.db.Table("talk_records")
	query.Joins("left join users on users.id = talk_records.user_id")

	ids := req.MessageIds
	if len(ids) > 3 {
		ids = ids[:3]
	}

	query.Where("talk_records.id in ?", ids)

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
			item["text"] = strutil.MtSubstr(text, 0, 30)
		case entity.MsgTypeCode:
			item["nickname"] = row.Nickname
			item["text"] = "【代码消息】"
		case entity.MsgTypeFile:
			item["nickname"] = row.Nickname
			item["text"] = "【文件消息】"
		}

		data = append(data, item)
	}

	return jsonutil.Encode(data), nil
}
