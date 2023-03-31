package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/timeutil"
)

type TalkMessageService struct {
	*repo.Source
	config              *config.Config
	unreadTalkCache     *cache.UnreadStorage
	lastMessage         *cache.MessageStorage
	talkRecordsVoteRepo *repo.TalkRecordsVote
	groupMemberRepo     *repo.GroupMember
	sidServer           *cache.ServerStorage
	client              *cache.ClientStorage
	fileSystem          *filesystem.Filesystem
	splitUploadDao      *repo.SplitUpload
	sequence            *repo.Sequence
}

func NewTalkMessageService(source *repo.Source, config *config.Config, unreadTalkCache *cache.UnreadStorage, lastMessage *cache.MessageStorage, talkRecordsVoteRepo *repo.TalkRecordsVote, groupMemberRepo *repo.GroupMember, sidServer *cache.ServerStorage, client *cache.ClientStorage, fileSystem *filesystem.Filesystem, splitUploadDao *repo.SplitUpload, sequence *repo.Sequence) *TalkMessageService {
	return &TalkMessageService{Source: source, config: config, unreadTalkCache: unreadTalkCache, lastMessage: lastMessage, talkRecordsVoteRepo: talkRecordsVoteRepo, groupMemberRepo: groupMemberRepo, sidServer: sidServer, client: client, fileSystem: fileSystem, splitUploadDao: splitUploadDao, sequence: sequence}
}

type SysTextMessageOpt struct {
	UserId     int
	TalkType   int
	ReceiverId int
	Text       string
}

// SendRevokeRecordMessage 撤销推送消息
func (s *TalkMessageService) SendRevokeRecordMessage(ctx context.Context, uid int, recordId int) error {
	var (
		err    error
		record model.TalkRecords
	)

	if err = s.Db().First(&record, recordId).Error; err != nil {
		return err
	}

	if record.IsRevoke == 1 {
		return nil
	}

	if record.UserId != uid {
		return errors.New("无权撤回回消息")
	}

	if time.Now().Unix() > record.CreatedAt.Add(3*time.Minute).Unix() {
		return errors.New("超出有效撤回时间范围，无法进行撤销！")
	}

	if err = s.Db().Model(&model.TalkRecords{Id: recordId}).Update("is_revoke", 1).Error; err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.EventTalkRevoke,
		"data": jsonutil.Encode(map[string]any{
			"record_id": record.Id,
		}),
	}

	s.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))

	return nil
}

type VoteMessageHandleOpt struct {
	UserId   int
	RecordId int
	Options  string
}

// VoteHandle 投票处理
func (s *TalkMessageService) VoteHandle(ctx context.Context, opts *VoteMessageHandleOpt) (int, error) {
	var (
		err  error
		vote *model.QueryVoteModel
	)

	tx := s.Db().Table("talk_records")
	tx.Select([]string{
		"talk_records.receiver_id", "talk_records.talk_type", "talk_records.msg_type",
		"vote.id as vote_id", "vote.id as record_id", "vote.answer_mode", "vote.answer_option",
		"vote.answer_num", "vote.status as vote_status",
	})
	tx.Joins("left join talk_records_vote as vote on vote.record_id = talk_records.id")
	tx.Where("talk_records.id = ?", opts.RecordId)

	res := tx.Take(&vote)
	if err := res.Error; err != nil {
		return 0, err
	}

	if res.RowsAffected == 0 {
		return 0, fmt.Errorf("投票信息不存在[%d]", opts.RecordId)
	}

	if vote.MsgType != entity.MsgTypeVote {
		return 0, fmt.Errorf("当前记录属于投票信息[%d]", vote.MsgType)
	}

	// 判断是否有投票权限

	var count int64
	s.Db().Table("talk_records_vote_answer").Where("vote_id = ? and user_id = ？", vote.VoteId, opts.UserId).Count(&count)
	if count > 0 { // 判断是否已投票
		return 0, fmt.Errorf("不能重复投票[%d]", vote.VoteId)
	}

	options := strings.Split(opts.Options, ",")
	sort.Strings(options)

	var answerOptions map[string]any
	if err = jsonutil.Decode(vote.AnswerOption, &answerOptions); err != nil {
		return 0, err
	}

	for _, option := range options {
		if _, ok := answerOptions[option]; !ok {
			return 0, fmt.Errorf("的投票选项不存在[%s]", option)
		}
	}

	// 判断是否单选
	if vote.AnswerMode == 0 {
		options = options[:1]
	}

	answers := make([]*model.TalkRecordsVoteAnswer, 0, len(options))

	for _, option := range options {
		answers = append(answers, &model.TalkRecordsVoteAnswer{
			VoteId: vote.VoteId,
			UserId: opts.UserId,
			Option: option,
		})
	}

	err = s.Db().Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("talk_records_vote").Where("id = ?", vote.VoteId).Updates(map[string]any{
			"answered_num": gorm.Expr("answered_num + 1"),
			"status":       gorm.Expr("if(answered_num >= answer_num, 1, 0)"),
		}).Error; err != nil {
			return err
		}

		if err = tx.Create(answers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	_, _ = s.talkRecordsVoteRepo.SetVoteAnswerUser(ctx, vote.VoteId)
	_, _ = s.talkRecordsVoteRepo.SetVoteStatistics(ctx, vote.VoteId)

	return vote.VoteId, nil
}

// 发送消息后置处理
func (s *TalkMessageService) afterHandle(ctx context.Context, record *model.TalkRecords, opts map[string]string) {

	if record.TalkType == entity.ChatPrivateMode {
		s.unreadTalkCache.Incr(ctx, entity.ChatPrivateMode, record.UserId, record.ReceiverId)

		if record.MsgType == entity.MsgTypeSystemText {
			s.unreadTalkCache.Incr(ctx, 1, record.ReceiverId, record.UserId)
		}
	} else if record.TalkType == entity.ChatGroupMode {

		// todo 需要加缓存
		ids := s.groupMemberRepo.GetMemberIds(ctx, record.ReceiverId)
		for _, uid := range ids {

			if uid == record.UserId {
				continue
			}

			s.unreadTalkCache.Incr(ctx, entity.ChatGroupMode, record.ReceiverId, uid)
		}
	}

	_ = s.lastMessage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  opts["text"],
		Datetime: timeutil.DateTime(),
	})

	content := jsonutil.Encode(map[string]any{
		"event": entity.EventTalk,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"talk_type":   record.TalkType,
			"record_id":   record.Id,
		}),
	})

	// 点对点消息采用精确投递
	if record.TalkType == entity.ChatPrivateMode {
		sids := s.sidServer.All(ctx, 1)

		// 小于三台服务器则采用全局广播
		if len(sids) <= 3 {
			s.Redis().Publish(ctx, entity.ImTopicChat, content)
		} else {
			for _, sid := range s.sidServer.All(ctx, 1) {
				for _, uid := range []int{record.UserId, record.ReceiverId} {
					if s.client.IsCurrentServerOnline(ctx, sid, entity.ImChannelChat, strconv.Itoa(uid)) {
						s.Redis().Publish(ctx, fmt.Sprintf(entity.ImTopicChatPrivate, sid), content)
					}
				}
			}
		}
	} else {
		s.Redis().Publish(ctx, entity.ImTopicChat, content)
	}
}
