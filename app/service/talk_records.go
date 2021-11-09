package service

import (
	"context"
	"go-chat/app/entity"
	"go-chat/app/http/dto"
	"go-chat/app/model"
	"go-chat/app/pkg/timeutil"
	"time"
)

type TalkRecordsService struct {
	*BaseService
}

type QueryTalkRecordsOpts struct {
	TalkType   int `json:"talk_type"`   // 对话类型
	UserId     int `json:"user_id"`     // 获取消息的用户
	ReceiverId int `json:"receiver_id"` // 接收者ID
	RecordId   int `json:"record_id"`   // 上次查询的最小消息ID
	Limit      int `json:"limit"`       // 数据行数
}

type QueryTalkRecordsItem struct {
	ID         int       `json:"id"`
	TalkType   int       `json:"talk_type"`
	MsgType    int       `json:"msg_type"`
	UserId     int       `json:"user_id"`
	ReceiverId int       `json:"receiver_id"`
	IsRevoke   int       `json:"is_revoke"`
	IsMark     int       `json:"is_mark"`
	IsRead     int       `json:"is_read"`
	QuoteId    int       `json:"quote_id"`
	WarnUsers  string    `json:"warn_users"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
}

func NewTalkRecordsService(base *BaseService) *TalkRecordsService {
	return &TalkRecordsService{base}
}

// GetTalkRecords 获取对话消息
func (s TalkRecordsService) GetTalkRecords(ctx context.Context, query *QueryTalkRecordsOpts) ([]*dto.TalkRecordsItem, error) {
	var (
		err    error
		items  []*QueryTalkRecordsItem
		fields = []string{
			"lar_talk_records.id",
			"lar_talk_records.talk_type",
			"lar_talk_records.msg_type",
			"lar_talk_records.user_id",
			"lar_talk_records.receiver_id",
			"lar_talk_records.is_revoke",
			"lar_talk_records.content",
			"lar_talk_records.created_at",
			"lar_users.nickname",
			"lar_users.avatar as avatar",
		}
	)

	tx := s.db.Debug().Table("lar_talk_records")

	tx.Joins("left join lar_users on lar_talk_records.user_id = lar_users.id")

	if query.RecordId > 0 {
		tx.Where("lar_talk_records.id < ?", query.RecordId)
	}

	if query.TalkType == entity.PrivateChat {
		subWhere := s.db.
			Where("lar_talk_records.user_id = ? and lar_talk_records.receiver_id = ?", query.UserId, query.ReceiverId).
			Or("lar_talk_records.user_id = ? and lar_talk_records.receiver_id = ?", query.ReceiverId, query.UserId)

		tx.Where(subWhere)
	} else {
		tx.Where("lar_talk_records.receiver_id = ?", query.ReceiverId)
	}

	tx.Where("lar_talk_records.talk_type = ?", query.TalkType)
	tx.Where("NOT EXISTS (SELECT 1 FROM `lar_talk_records_delete` WHERE lar_talk_records_delete.record_id = lar_talk_records.id AND lar_talk_records_delete.user_id = ? LIMIT 1)", query.UserId)
	tx.Select(fields).Order("lar_talk_records.id desc").Limit(query.Limit)

	if err = tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return s.HandleTalkRecords(items)
}

// SearchTalkRecords 对话搜索消息
func (s *TalkRecordsService) SearchTalkRecords() {

}

func (s *TalkRecordsService) HandleTalkRecords(items []*QueryTalkRecordsItem) ([]*dto.TalkRecordsItem, error) {
	var (
		files     []int
		codes     []int
		forwards  []int
		invites   []int
		votes     []int
		logins    []int
		locations []int

		fileItems     []*model.TalkRecordsFile
		codeItems     []*model.TalkRecordsCode
		forwardItems  []*model.TalkRecordsForward
		inviteItems   []*model.TalkRecordsInvite
		voteItems     []*model.TalkRecordsVote
		loginItems    []*model.TalkRecordsLogin
		locationItems []*model.TalkRecordsLocation
	)

	for _, item := range items {
		switch item.MsgType {
		case entity.MsgTypeFile:
			files = append(files, item.ID)
		case entity.MsgTypeForward:
			forwards = append(forwards, item.ID)
		case entity.MsgTypeCode:
			codes = append(codes, item.ID)
		case entity.MsgTypeVote:
			votes = append(votes, item.ID)
		case entity.MsgTypeGroupNotice:
		case entity.MsgTypeFriendApply:
		case entity.MsgTypeUserLogin:
			logins = append(logins, item.ID)
		case entity.MsgTypeGroupInvite:
			invites = append(invites, item.ID)
		case entity.MsgTypeLocation:
			locations = append(locations, item.ID)
		}
	}

	if len(files) > 0 {
		s.db.Model(model.TalkRecordsFile{}).Where("record_id in ?", files).Scan(&fileItems)
	}

	if len(forwards) > 0 {
		s.db.Model(model.TalkRecordsForward{}).Where("record_id in ?", forwards).Scan(&forwardItems)
	}

	if len(codes) > 0 {
		s.db.Model(model.TalkRecordsCode{}).Where("record_id in ?", codes).Scan(&codeItems)
	}

	if len(votes) > 0 {
		s.db.Model(model.TalkRecordsCode{}).Where("record_id in ?", votes).Scan(&voteItems)
	}

	if len(logins) > 0 {
		s.db.Model(model.TalkRecordsLogin{}).Where("record_id in ?", votes).Scan(&loginItems)
	}

	if len(invites) > 0 {
		s.db.Model(model.TalkRecordsInvite{}).Where("record_id in ?", invites).Scan(&inviteItems)
	}

	if len(locations) > 0 {
		s.db.Model(model.TalkRecordsLocation{}).Where("record_id in ?", locations).Scan(&locationItems)
	}

	newItems := make([]*dto.TalkRecordsItem, 0, len(items))

	for _, item := range items {
		data := &dto.TalkRecordsItem{
			ID:         item.ID,
			TalkType:   item.TalkType,
			MsgType:    item.MsgType,
			UserID:     item.UserId,
			ReceiverID: item.ReceiverId,
			Nickname:   item.Nickname,
			Avatar:     item.Avatar,
			IsRevoke:   item.IsRevoke,
			IsMark:     item.IsMark,
			IsRead:     item.IsRead,
			Content:    item.Content,
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
		}

		newItems = append(newItems, data)
	}

	return newItems, nil
}
