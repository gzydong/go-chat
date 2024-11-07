package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ ITalkRecordService = (*TalkRecordService)(nil)

type ITalkRecordService interface {
	FindTalkPrivateRecord(ctx context.Context, uid int, msgId string) (*TalkRecord, error)
	FindTalkGroupRecord(ctx context.Context, msgId string) (*TalkRecord, error)
	FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*TalkRecord, error)
	FindForwardRecords(ctx context.Context, uid int, msgIds []string, talkType int) ([]*TalkRecord, error)
}

type TalkRecord struct {
	MsgId     string `json:"msg_id"`
	Sequence  int    `json:"sequence"`
	MsgType   int    `json:"msg_type"`
	UserId    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsRevoked int    `json:"is_revoked"`
	SendTime  string `json:"send_time"`
	Extra     any    `json:"extra"` // 额外参数
	Quote     any    `json:"quote"` // 额外参数
}

type TalkRecordService struct {
	*repo.Source
	TalkVoteCache         *cache.Vote
	TalkRecordsVoteRepo   *repo.GroupVote
	GroupMemberRepo       *repo.GroupMember
	TalkRecordFriendRepo  *repo.TalkUserMessage
	TalkRecordGroupRepo   *repo.TalkGroupMessage
	TalkRecordsDeleteRepo *repo.TalkGroupMessageDel
}

type FindAllTalkRecordsOpt struct {
	TalkType   int   // 对话类型
	UserId     int   // 获取消息的用户
	ReceiverId int   // 接收者ID
	MsgType    []int // 消息类型
	Cursor     int   // 上次查询的游标
	Limit      int   // 数据行数
}

type QueryTalkRecord struct {
	MsgId     string    `json:"msg_id"`
	Sequence  int64     `json:"sequence"`
	MsgType   int       `json:"msg_type"`
	UserId    int       `json:"user_id"`
	IsRevoked int       `json:"is_revoked"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Extra     string    `json:"extra"`
	Quote     string    `json:"quote"`
	SendTime  time.Time `json:"send_time"`
}

func (s *TalkRecordService) FindTalkPrivateRecord(ctx context.Context, uid int, msgId string) (*TalkRecord, error) {
	talkRecordFriendInfo, err := s.TalkRecordFriendRepo.FindByWhere(ctx, "msg_id = ? and user_id = ?", msgId, uid)
	if err != nil {
		return nil, err
	}

	record := &QueryTalkRecord{
		MsgId:     talkRecordFriendInfo.MsgId,
		Sequence:  talkRecordFriendInfo.Sequence,
		MsgType:   talkRecordFriendInfo.MsgType,
		UserId:    talkRecordFriendInfo.FromId,
		IsRevoked: talkRecordFriendInfo.IsRevoked,
		Nickname:  "",
		Avatar:    "",
		Extra:     talkRecordFriendInfo.Extra,
		Quote:     talkRecordFriendInfo.Quote,
		SendTime:  talkRecordFriendInfo.SendTime,
	}

	list, err := s.handleTalkRecords(ctx, []*QueryTalkRecord{record})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

func (s *TalkRecordService) FindTalkGroupRecord(ctx context.Context, msgId string) (*TalkRecord, error) {
	talkRecordGroupInfo, err := s.TalkRecordGroupRepo.FindByMsgId(ctx, msgId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if talkRecordGroupInfo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	record := &QueryTalkRecord{
		MsgId:     talkRecordGroupInfo.MsgId,
		Sequence:  talkRecordGroupInfo.Sequence,
		MsgType:   talkRecordGroupInfo.MsgType,
		UserId:    talkRecordGroupInfo.FromId,
		IsRevoked: talkRecordGroupInfo.IsRevoked,
		Nickname:  "",
		Avatar:    "",
		Extra:     talkRecordGroupInfo.Extra,
		Quote:     talkRecordGroupInfo.Quote,
		SendTime:  talkRecordGroupInfo.SendTime,
	}

	list, err := s.handleTalkRecords(ctx, []*QueryTalkRecord{record})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

// FindAllTalkRecords 获取所有对话消息
func (s *TalkRecordService) FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*TalkRecord, error) {
	var (
		items  = make([]*QueryTalkRecord, 0, opt.Limit)
		cursor = opt.Cursor
	)

	for {
		// 这里查询数据放弃了关联查询，所以这里需要查询多次，防止查询中存在用户已删除的数据需要过滤掉
		list, err := s.findAllRecords(ctx, &FindAllTalkRecordsOpt{
			TalkType:   opt.TalkType,
			UserId:     opt.UserId,
			ReceiverId: opt.ReceiverId,
			MsgType:    opt.MsgType,
			Cursor:     cursor,
			Limit:      opt.Limit + 10, // 多查几条数据
		})

		if err != nil {
			return nil, err
		}

		if len(list) == 0 {
			break
		}

		if opt.TalkType == entity.ChatGroupMode {
			tmpMsgIds := make([]string, 0, len(list))
			for _, v := range list {
				tmpMsgIds = append(tmpMsgIds, v.MsgId)
			}

			msgIds, err := s.TalkRecordsDeleteRepo.FindAllMsgIds(ctx, opt.UserId, tmpMsgIds)
			if err != nil {
				return nil, err
			}

			hashIds := make(map[string]struct{}, len(msgIds))
			for _, msgId := range msgIds {
				hashIds[msgId] = struct{}{}
			}

			for _, v := range list {
				if _, ok := hashIds[v.MsgId]; ok {
					continue
				}

				items = append(items, v)
			}
		} else {
			items = append(items, list...)
		}

		if len(items) >= opt.Limit || len(list) < opt.Limit {
			break
		}

		// 设置游标继续往下执行
		cursor = int(list[len(list)-1].Sequence)
	}

	if len(items) > opt.Limit {
		items = items[:opt.Limit]
	}

	return s.handleTalkRecords(ctx, items)
}

func (s *TalkRecordService) findAllRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*QueryTalkRecord, error) {
	query := s.Source.Db().WithContext(ctx)

	fields := []string{
		"msg_id",
		"sequence",
		"msg_type",
		"is_revoked",
		"extra",
		"quote",
		"send_time",
		"from_id as user_id",
	}

	if opt.TalkType == 1 {
		query = query.Table("talk_user_message")
		query.Where("user_id = ?", opt.UserId)
		query.Where("to_from_id = ?", opt.ReceiverId)
		query.Where("is_deleted = ?", model.No)
	} else {
		query = query.Table("talk_group_message")
		query.Where("group_id = ?", opt.ReceiverId)
	}

	query.Select(fields)

	if opt.Cursor > 0 {
		query.Where("sequence < ?", opt.Cursor)
	}

	if len(opt.MsgType) > 0 {
		query.Where("msg_type in ?", opt.MsgType)
	}

	query.Order("sequence desc").Limit(opt.Limit)

	var items []*QueryTalkRecord
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindForwardRecords 获取转发消息记录
func (s *TalkRecordService) FindForwardRecords(ctx context.Context, uid int, msgIds []string, talkType int) ([]*TalkRecord, error) {
	var (
		fields = []string{
			"msg_id",
			"sequence",
			"msg_type",
			"is_revoked",
			"extra",
			"quote",
			"send_time",
			"from_id as user_id",
		}
		items     = make([]*QueryTalkRecord, 0)
		tableName = "talk_user_message"
	)

	if talkType == 1 {

	} else {
		tableName = "talk_group_message"
	}

	query := s.Source.Db().Table(tableName)
	query.Select(fields)
	query.Where("msg_id in ?", msgIds)
	query.Order("sequence asc")

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return s.handleTalkRecords(ctx, items)
}

// HandleTalkRecords 处理消息
func (s *TalkRecordService) handleTalkRecords(ctx context.Context, items []*QueryTalkRecord) ([]*TalkRecord, error) {
	if len(items) == 0 {
		return make([]*TalkRecord, 0), nil
	}

	uids := make([]int, 0, len(items))
	for _, item := range items {
		uids = append(uids, item.UserId)
	}

	var usersItems []*model.Users
	err := s.Source.Db().Model(&model.Users{}).Select("id,nickname,avatar").Where("id in ?", sliceutil.Unique(uids)).Scan(&usersItems).Error
	if err != nil {
		return nil, err
	}

	hashUser := make(map[int]*model.Users)
	for _, user := range usersItems {
		hashUser[user.Id] = user
	}

	newItems := make([]*TalkRecord, 0, len(items))
	for _, item := range items {
		data := &TalkRecord{
			MsgId:     item.MsgId,
			Sequence:  int(item.Sequence),
			MsgType:   item.MsgType,
			UserId:    item.UserId,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			IsRevoked: item.IsRevoked,
			SendTime:  item.SendTime.Format(time.DateTime),
			Extra:     make(map[string]any),
			Quote:     make(map[string]any),
		}

		if user, ok := hashUser[item.UserId]; ok {
			data.Nickname = user.Nickname
			data.Avatar = user.Avatar
		}

		if err := jsonutil.Decode(item.Extra, &data.Extra); err != nil {
			fmt.Println("ERR===>", item.MsgId, data.Extra, item.Extra)
		}

		if err := jsonutil.Decode(item.Quote, &data.Quote); err != nil {
			fmt.Println("ERR===>", item.MsgId, data.Quote, item.Quote)
		}

		newItems = append(newItems, data)
	}

	return newItems, nil
}
