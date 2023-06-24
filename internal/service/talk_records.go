package service

import (
	"context"
	"sort"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type TalkRecordsItem struct {
	Id         int    `json:"id"`
	Sequence   int    `json:"sequence"`
	MsgId      string `json:"msg_id"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int    `json:"user_id"`
	ReceiverId int    `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	Extra      any    `json:"extra"` // 额外参数
}

type TalkRecordsService struct {
	*repo.Source
	talkVoteCache       *cache.Vote
	talkRecordsVoteRepo *repo.TalkRecordsVote
	groupMemberRepo     *repo.GroupMember
	talkRecordsRepo     *repo.TalkRecords
}

func NewTalkRecordsService(source *repo.Source, talkVoteCache *cache.Vote, talkRecordsVoteRepo *repo.TalkRecordsVote, groupMemberRepo *repo.GroupMember, repo *repo.TalkRecords) *TalkRecordsService {
	return &TalkRecordsService{Source: source, talkVoteCache: talkVoteCache, talkRecordsVoteRepo: talkRecordsVoteRepo, groupMemberRepo: groupMemberRepo, talkRecordsRepo: repo}
}

func (s *TalkRecordsService) Dao() *repo.TalkRecords {
	return s.talkRecordsRepo
}

type QueryTalkRecordsOpt struct {
	TalkType   int   // 对话类型
	UserId     int   // 获取消息的用户
	ReceiverId int   // 接收者ID
	MsgType    []int // 消息类型
	RecordId   int   // 上次查询的最小消息ID
	Limit      int   // 数据行数
}

type QueryTalkRecordsItem struct {
	Id         int       `json:"id"`
	MsgId      string    `json:"msg_id"`
	Sequence   int64     `json:"sequence"`
	TalkType   int       `json:"talk_type"`
	MsgType    int       `json:"msg_type"`
	UserId     int       `json:"user_id"`
	ReceiverId int       `json:"receiver_id"`
	IsRevoke   int       `json:"is_revoke"`
	IsMark     int       `json:"is_mark"`
	IsRead     int       `json:"is_read"`
	QuoteId    int       `json:"quote_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Extra      string    `json:"extra"`
}

// GetTalkRecords 获取对话消息
func (s *TalkRecordsService) GetTalkRecords(ctx context.Context, opt *QueryTalkRecordsOpt) ([]*TalkRecordsItem, error) {
	var (
		items  = make([]*QueryTalkRecordsItem, 0, opt.Limit)
		fields = []string{
			"talk_records.id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.msg_id",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.is_read",
			"talk_records.content",
			"talk_records.extra",
			"talk_records.created_at",
			"users.nickname",
			"users.avatar as avatar",
		}
	)

	query := s.Db().WithContext(ctx).Table("talk_records")
	query.Joins("left join users on talk_records.user_id = users.id")
	query.Joins("left join talk_records_delete on talk_records.id = talk_records_delete.record_id and talk_records_delete.user_id = ?", opt.UserId)

	if opt.RecordId > 0 {
		query.Where("talk_records.sequence < ?", opt.RecordId)
	}

	if opt.TalkType == entity.ChatPrivateMode {
		subQuery := s.Db().Where("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.UserId, opt.ReceiverId)
		subQuery.Or("talk_records.user_id = ? and talk_records.receiver_id = ?", opt.ReceiverId, opt.UserId)

		query.Where(subQuery)
	} else {
		query.Where("talk_records.receiver_id = ?", opt.ReceiverId)
	}

	if opt.MsgType != nil && len(opt.MsgType) > 0 {
		query.Where("talk_records.msg_type in ?", opt.MsgType)
	}

	query.Where("talk_records.talk_type = ?", opt.TalkType)
	query.Where("ifnull(talk_records_delete.id,0) = 0")
	query.Select(fields).Order("talk_records.sequence desc").Limit(opt.Limit)

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return make([]*TalkRecordsItem, 0), nil
	}

	return s.HandleTalkRecords(ctx, items)
}

// SearchTalkRecords 对话搜索消息
func (s *TalkRecordsService) SearchTalkRecords() {

}

func (s *TalkRecordsService) GetTalkRecord(ctx context.Context, recordId int64) (*TalkRecordsItem, error) {
	var (
		err    error
		item   *QueryTalkRecordsItem
		fields = []string{
			"talk_records.id",
			"talk_records.msg_id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.content",
			"talk_records.extra",
			"talk_records.created_at",
			"users.nickname",
			"users.avatar as avatar",
		}
	)

	query := s.Db().Table("talk_records")
	query.Joins("left join users on talk_records.user_id = users.id")
	query.Where("talk_records.id = ?", recordId)

	if err = query.Select(fields).Take(&item).Error; err != nil {
		return nil, err
	}

	list, err := s.HandleTalkRecords(ctx, []*QueryTalkRecordsItem{item})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

// GetForwardRecords 获取转发消息记录
func (s *TalkRecordsService) GetForwardRecords(ctx context.Context, uid int, recordId int64) ([]*TalkRecordsItem, error) {

	record, err := s.talkRecordsRepo.FindById(ctx, int(recordId))
	if err != nil {
		return nil, err
	}

	// if record.TalkType == entity.ChatPrivateMode {
	// 	if record.UserId != uid && record.ReceiverId != uid {
	// 		return nil, entity.ErrPermissionDenied
	// 	}
	// } else if record.TalkType == entity.ChatGroupMode {
	// 	if !s.groupMemberRepo.IsMember(ctx, record.ReceiverId, uid, true) {
	// 		return nil, entity.ErrPermissionDenied
	// 	}
	// } else {
	// 	return nil, entity.ErrPermissionDenied
	// }

	var extra model.TalkRecordExtraForward
	if err := jsonutil.Decode(record.Extra, &extra); err != nil {
		return nil, err
	}

	var (
		items  = make([]*QueryTalkRecordsItem, 0)
		fields = []string{
			"talk_records.id",
			"talk_records.msg_id",
			"talk_records.sequence",
			"talk_records.talk_type",
			"talk_records.msg_type",
			"talk_records.user_id",
			"talk_records.receiver_id",
			"talk_records.is_revoke",
			"talk_records.content",
			"talk_records.extra",
			"talk_records.created_at",
			"users.nickname",
			"users.avatar as avatar",
		}
	)

	query := s.Db().Table("talk_records")
	query.Select(fields)
	query.Joins("left join users on talk_records.user_id = users.id")
	query.Where("talk_records.id in ?", extra.MsgIds)
	query.Order("talk_records.sequence asc")

	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return s.HandleTalkRecords(ctx, items)
}

func (s *TalkRecordsService) HandleTalkRecords(ctx context.Context, items []*QueryTalkRecordsItem) ([]*TalkRecordsItem, error) {
	var (
		votes     []int
		voteItems []*model.TalkRecordsVote
	)

	for _, item := range items {
		switch item.MsgType {
		case entity.ChatMsgTypeVote:
			votes = append(votes, item.Id)
		}
	}

	hashVotes := make(map[int]*model.TalkRecordsVote)
	if len(votes) > 0 {
		s.Db().Model(&model.TalkRecordsVote{}).Where("record_id in ?", votes).Scan(&voteItems)
		for i := range voteItems {
			hashVotes[voteItems[i].RecordId] = voteItems[i]
		}
	}

	newItems := make([]*TalkRecordsItem, 0, len(items))
	for _, item := range items {
		data := &TalkRecordsItem{
			Id:         item.Id,
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
			IsRead:     item.IsRead,
			Content:    item.Content,
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
			Extra:      make(map[string]any),
		}

		_ = jsonutil.Decode(item.Extra, &data.Extra)

		switch item.MsgType {
		case entity.ChatMsgTypeVote:
			if value, ok := hashVotes[item.Id]; ok {
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
				if uids, err := s.talkRecordsVoteRepo.GetVoteAnswerUser(ctx, value.Id); err == nil {
					users = uids
				}

				var statistics any

				if res, err := s.talkRecordsVoteRepo.GetVoteStatistics(ctx, value.Id); err != nil {
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
						"record_id":     value.RecordId,
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
