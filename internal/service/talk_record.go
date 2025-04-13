package service

import (
	"context"
	"errors"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ ITalkRecordService = (*TalkRecordService)(nil)

type ITalkRecordService interface {
	FindAllPrivateRecordByOriMsgId(ctx context.Context, msgId string) ([]*model.TalkUserMessage, error)
	FindPrivateRecordByMsgId(ctx context.Context, msgId string) (*model.TalkUserMessage, error)
	FindTalkPrivateRecord(ctx context.Context, uid int, msgId string) (*model.TalkMessageRecord, error)
	FindTalkGroupRecord(ctx context.Context, msgId string) (*model.TalkMessageRecord, error)
	FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*model.TalkMessageRecord, error)
	FindForwardRecords(ctx context.Context, uid int, msgIds []string, talkType int) ([]*model.TalkMessageRecord, error)
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

func (s *TalkRecordService) FindPrivateRecordByMsgId(ctx context.Context, msgId string) (*model.TalkUserMessage, error) {
	return s.TalkRecordFriendRepo.FindByMsgId(ctx, msgId)
}

func (s *TalkRecordService) FindAllPrivateRecordByOriMsgId(ctx context.Context, msgId string) ([]*model.TalkUserMessage, error) {
	return s.TalkRecordFriendRepo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("org_msg_id = ?", msgId)
	})
}

type FindAllTalkRecordsOpt struct {
	TalkType   int   // 对话类型
	UserId     int   // 获取消息的用户
	ReceiverId int   // 接收者ID
	MsgType    []int // 消息类型
	Cursor     int   // 上次查询的游标
	Limit      int   // 数据行数
}

func (s *TalkRecordService) FindTalkPrivateRecord(ctx context.Context, uid int, msgId string) (*model.TalkMessageRecord, error) {
	talkRecordFriendInfo, err := s.TalkRecordFriendRepo.FindByWhere(ctx, "msg_id = ? and user_id = ?", msgId, uid)
	if err != nil {
		return nil, err
	}

	record := &model.TalkMessageRecord{
		TalkMode:  entity.ChatPrivateMode,
		FromId:    talkRecordFriendInfo.FromId,
		ToFromId:  talkRecordFriendInfo.ToFromId,
		MsgId:     talkRecordFriendInfo.MsgId,
		Sequence:  int(talkRecordFriendInfo.Sequence),
		MsgType:   talkRecordFriendInfo.MsgType,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: talkRecordFriendInfo.IsRevoked,
		SendTime:  talkRecordFriendInfo.SendTime,
		Extra:     talkRecordFriendInfo.Extra,
		Quote:     talkRecordFriendInfo.Quote,
	}

	list, err := s.handleTalkRecords(ctx, []*model.TalkMessageRecord{record})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

func (s *TalkRecordService) FindTalkGroupRecord(ctx context.Context, msgId string) (*model.TalkMessageRecord, error) {
	talkRecordGroupInfo, err := s.TalkRecordGroupRepo.FindByMsgId(ctx, msgId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if talkRecordGroupInfo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	record := &model.TalkMessageRecord{
		TalkMode:  entity.ChatGroupMode,
		FromId:    talkRecordGroupInfo.FromId,
		ToFromId:  talkRecordGroupInfo.GroupId,
		MsgId:     talkRecordGroupInfo.MsgId,
		Sequence:  int(talkRecordGroupInfo.Sequence),
		MsgType:   talkRecordGroupInfo.MsgType,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: talkRecordGroupInfo.IsRevoked,
		SendTime:  talkRecordGroupInfo.SendTime,
		Extra:     talkRecordGroupInfo.Extra,
		Quote:     talkRecordGroupInfo.Quote,
	}

	list, err := s.handleTalkRecords(ctx, []*model.TalkMessageRecord{record})
	if err != nil {
		return nil, err
	}

	return list[0], nil
}

// FindAllTalkRecords 获取所有对话消息
func (s *TalkRecordService) FindAllTalkRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*model.TalkMessageRecord, error) {
	var (
		items  = make([]*model.TalkMessageRecord, 0, opt.Limit)
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
		cursor = list[len(list)-1].Sequence
	}

	if len(items) > opt.Limit {
		items = items[:opt.Limit]
	}

	return s.handleTalkRecords(ctx, items)
}

func (s *TalkRecordService) findAllRecords(ctx context.Context, opt *FindAllTalkRecordsOpt) ([]*model.TalkMessageRecord, error) {
	query := s.Source.Db().WithContext(ctx)

	fields := []string{
		"msg_id",
		"sequence",
		"msg_type",
		"is_revoked",
		"extra",
		"quote",
		"send_time",
		"from_id",
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

	var items []*model.TalkMessageRecord
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(items); i++ {
		items[i].TalkMode = opt.TalkType
		items[i].ToFromId = opt.ReceiverId
	}

	return items, nil
}

// FindForwardRecords 获取转发消息记录
func (s *TalkRecordService) FindForwardRecords(ctx context.Context, uid int, msgIds []string, talkType int) ([]*model.TalkMessageRecord, error) {
	var (
		fields = []string{
			"msg_id",
			"sequence",
			"msg_type",
			"is_revoked",
			"extra",
			"quote",
			"send_time",
			"from_id",
		}
		items     = make([]*model.TalkMessageRecord, 0)
		tableName = "talk_user_message"
	)

	if talkType == 2 {
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
func (s *TalkRecordService) handleTalkRecords(ctx context.Context, items []*model.TalkMessageRecord) ([]*model.TalkMessageRecord, error) {
	if len(items) == 0 {
		return make([]*model.TalkMessageRecord, 0), nil
	}

	uids := make([]int, 0, len(items))
	for _, item := range items {
		uids = append(uids, item.FromId)
	}

	var usersItems []*model.Users
	err := s.Source.Db().WithContext(ctx).Model(&model.Users{}).
		Select("id,nickname,avatar").
		Where("id in ?", sliceutil.Unique(uids)).Scan(&usersItems).Error
	if err != nil {
		return nil, err
	}

	hashUser := make(map[int]*model.Users)
	for _, user := range usersItems {
		hashUser[user.Id] = user
	}

	for i := 0; i < len(items); i++ {
		if user, ok := hashUser[items[i].FromId]; ok {
			items[i].Nickname = user.Nickname
			items[i].Avatar = user.Avatar
		}

		//if err = jsonutil.Unmarshal(items[i].Extra, &items[i].Extra); err != nil {
		//	fmt.Println("ERR===>", items[i].MsgId, items[i].Extra, items[i].Extra)
		//}
		//
		//if err = jsonutil.Unmarshal(items[i].Quote, &items[i].Quote); err != nil {
		//	fmt.Println("ERR===>", items[i].MsgId, items[i].Quote, items[i].Quote)
		//}
	}

	return items, nil
}
