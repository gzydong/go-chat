package service

import (
	"context"
	"sort"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

var _ ITalkRecordsService = (*TalkRecordsService)(nil)

type ITalkRecordsService interface {
	FindTalkRecord(ctx context.Context, msgId string) (*TalkRecord, error)
	FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*TalkRecord, error)
	FindForwardRecords(ctx context.Context, uid int, msgId string) ([]*TalkRecord, error)
}

type TalkRecord struct {
	MsgId      string `json:"msg_id"`
	Sequence   int    `json:"sequence"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int    `json:"user_id"`
	ReceiverId int    `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	CreatedAt  string `json:"created_at"`
	Extra      any    `json:"extra"` // 额外参数
}

type TalkRecordsService struct {
	*repo.Source
	TalkVoteCache         *cache.Vote
	TalkRecordsVoteRepo   *repo.TalkRecordsVote
	GroupMemberRepo       *repo.GroupMember
	TalkRecordsRepo       *repo.TalkRecords
	TalkRecordsDeleteRepo *repo.TalkRecordsDelete
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
	MsgId      string    `json:"msg_id"`
	Sequence   int64     `json:"sequence"`
	TalkType   int       `json:"talk_type"`
	MsgType    int       `json:"msg_type"`
	UserId     int       `json:"user_id"`
	ReceiverId int       `json:"receiver_id"`
	IsRevoke   int       `json:"is_revoke"`
	IsMark     int       `json:"is_mark"`
	QuoteId    int       `json:"quote_id"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Extra      string    `json:"extra"`
	CreatedAt  time.Time `json:"created_at"`
}

// FindTalkRecord 获取对话消息
func (s *TalkRecordsService) FindTalkRecord(ctx context.Context, msgId string) (*TalkRecord, error) {
	var (
		err    error
		item   *QueryTalkRecord
		fields = []string{
			"talk_records.msg_id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.extra",
			"talk_records.created_at",
		}
	)

	query := s.Source.Db().Table("talk_records")
	query.Where("talk_records.msg_id = ?", msgId)

	if err = query.Select(fields).Take(&item).Error; err != nil {
		return nil, err
	}

	list, err := s.handleTalkRecords(ctx, []*QueryTalkRecord{item})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

// FindAllTalkRecords 获取所有对话消息
func (s *TalkRecordsService) FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*TalkRecord, error) {
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

		tmpMsgIds := make([]string, 0, len(list))
		for _, v := range list {
			tmpMsgIds = append(tmpMsgIds, v.MsgId)
		}

		msgIds, err := s.TalkRecordsDeleteRepo.FindAllMsgIds(ctx, tmpMsgIds, opt.UserId)
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

func (s *TalkRecordsService) findAllRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*QueryTalkRecord, error) {
	query := s.Source.Db().WithContext(ctx).Table("talk_records")
	query.Select([]string{
		"talk_records.sequence",
		"talk_records.talk_type",
		"talk_records.msg_type",
		"talk_records.msg_id",
		"talk_records.user_id",
		"talk_records.receiver_id",
		"talk_records.is_revoke",
		"talk_records.extra",
		"talk_records.created_at",
	})

	if opt.Cursor > 0 {
		query.Where("talk_records.sequence < ?", opt.Cursor)
	}

	if opt.TalkType == entity.ChatPrivateMode {
		subQuery := s.Source.Db().Where("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.UserId, opt.ReceiverId)
		subQuery.Or("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.ReceiverId, opt.UserId)

		query.Where(subQuery)
	} else {
		query.Where("talk_records.receiver_id = ?", opt.ReceiverId)
	}

	if opt.MsgType != nil && len(opt.MsgType) > 0 {
		query.Where("talk_records.msg_type in ?", opt.MsgType)
	}

	query.Where("talk_records.talk_type = ?", opt.TalkType)
	query.Order("talk_records.sequence desc").Limit(opt.Limit)

	var items []*QueryTalkRecord
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindForwardRecords 获取转发消息记录
func (s *TalkRecordsService) FindForwardRecords(ctx context.Context, uid int, msgId string) ([]*TalkRecord, error) {
	record, err := s.TalkRecordsRepo.FindByMsgId(ctx, msgId)
	if err != nil {
		return nil, err
	}

	var extra model.TalkRecordExtraForward
	if err := jsonutil.Decode(record.Extra, &extra); err != nil {
		return nil, err
	}

	var (
		items  = make([]*QueryTalkRecord, 0)
		fields = []string{
			"talk_records.msg_id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.extra",
			"talk_records.created_at",
		}
	)

	query := s.Source.Db().Table("talk_records")
	query.Select(fields)
	query.Where("talk_records.msg_id in ?", extra.MsgIds)
	query.Order("talk_records.sequence asc")

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return s.handleTalkRecords(ctx, items)
}

// HandleTalkRecords 处理消息
func (s *TalkRecordsService) handleTalkRecords(ctx context.Context, items []*QueryTalkRecord) ([]*TalkRecord, error) {
	if len(items) == 0 {
		return make([]*TalkRecord, 0), nil
	}

	var (
		votes     []string
		voteItems []*model.TalkRecordsVote
	)

	uids := make([]int, 0, len(items))
	for _, item := range items {
		uids = append(uids, item.UserId)

		switch item.MsgType {
		case entity.ChatMsgTypeVote:
			votes = append(votes, item.MsgId)
		}
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

	hashVotes := make(map[string]*model.TalkRecordsVote)
	if len(votes) > 0 {
		s.Source.Db().Model(&model.TalkRecordsVote{}).Where("msg_id in ?", votes).Scan(&voteItems)
		for i := range voteItems {
			hashVotes[voteItems[i].MsgId] = voteItems[i]
		}
	}

	newItems := make([]*TalkRecord, 0, len(items))
	for _, item := range items {
		data := &TalkRecord{
			MsgId:      item.MsgId,
			Sequence:   int(item.Sequence),
			TalkType:   item.TalkType,
			MsgType:    item.MsgType,
			UserId:     item.UserId,
			ReceiverId: item.ReceiverId,
			Nickname:   item.Nickname,
			Avatar:     item.Avatar,
			IsRevoke:   item.IsRevoke,
			IsMark:     item.IsMark,
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
			Extra:      make(map[string]any),
		}

		if user, ok := hashUser[item.UserId]; ok {
			data.Nickname = user.Nickname
			data.Avatar = user.Avatar
		}

		_ = jsonutil.Decode(item.Extra, &data.Extra)

		switch item.MsgType {
		case entity.ChatMsgTypeVote:
			if value, ok := hashVotes[item.MsgId]; ok {
				options := make(map[string]any)
				opts := make([]any, 0)

				if err := jsonutil.Decode(value.AnswerOption, &options); err == nil {
					arr := make([]string, 0, len(options))
					for k := range options {
						arr = append(arr, k)
					}

					sort.Strings(arr)

					for _, v := range arr {
						opts = append(opts, map[string]any{
							"key":   v,
							"value": options[v],
						})
					}
				}

				users := make([]int, 0)
				if uids, err := s.TalkRecordsVoteRepo.GetVoteAnswerUser(ctx, value.Id); err == nil {
					users = uids
				}

				var statistics any

				if res, err := s.TalkRecordsVoteRepo.GetVoteStatistics(ctx, value.Id); err != nil {
					statistics = map[string]any{
						"count":   0,
						"options": map[string]int{},
					}
				} else {
					statistics = res
				}

				data.Extra = map[string]any{
					"detail": map[string]any{
						"id":            value.Id,
						"msg_id":        value.MsgId,
						"title":         value.Title,
						"answer_mode":   value.AnswerMode,
						"status":        value.Status,
						"answer_option": opts,
						"answer_num":    value.AnswerNum,
						"answered_num":  value.AnsweredNum,
					},
					"statistics": statistics,
					"vote_users": users, // 已投票成员
				}
			}
		}

		newItems = append(newItems, data)
	}

	return newItems, nil
}
