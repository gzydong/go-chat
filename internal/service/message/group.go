package message

import (
	"context"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
)

func (s *Service) CreateGroupMessage(ctx context.Context, option CreateGroupMessageOption) error {
	quoteJsonText := "{}"

	if option.QuoteId != "" {
		quoteRecord := &model.TalkGroupMessage{}
		if err := s.Db().First(quoteRecord, "msg_id = ?", option.QuoteId).Error; err != nil {
			return err
		}

		user := &model.Users{}
		if err := s.Db().First(user, "id = ?", quoteRecord.FromId).Error; err != nil {
			return err
		}

		quote := &model.Quote{
			QuoteId: option.QuoteId,
			MsgType: 1,
		}

		quote.Nickname = user.Nickname
		quote.Content = s.getTextMessage(quoteRecord.MsgType, quoteRecord.Extra)
		quoteJsonText = jsonutil.Encode(quote)
	}

	item := &model.TalkGroupMessage{
		MsgId:     strutil.NewMsgId(),
		Sequence:  s.Sequence.Get(ctx, option.ToFromId, false),
		MsgType:   option.MsgType,
		GroupId:   option.ToFromId,
		FromId:    option.FromId,
		Quote:     quoteJsonText,
		Extra:     option.Extra,
		IsRevoked: model.No,
		SendTime:  time.Now(),
	}

	if err := s.Db().WithContext(ctx).Create(item).Error; err != nil {
		return err
	}

	err := s.PushMessage.Push(ctx, entity.ImTopicChat, &entity.SubscribeMessage{
		Event: entity.SubEventImMessage,
		Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
			TalkMode: entity.ChatGroupMode,
			Message:  jsonutil.Encode(item),
		}),
	})
	if err != nil {
		logger.Errorf("CreateGroupMessage publish message error:%s", err.Error())
	}

	pipe := s.Source.Redis().Pipeline()
	for _, uid := range s.GroupMemberRepo.GetMemberIds(ctx, item.GroupId) {
		if uid != item.FromId {
			s.UnreadStorage.PipeIncr(ctx, pipe, uid, entity.ChatGroupMode, item.GroupId)
		}
	}
	_, _ = pipe.Exec(ctx)

	// 更新最后一条消息
	_ = s.MessageStorage.Set(ctx, entity.ChatGroupMode, item.FromId, item.GroupId, &cache.LastCacheMessage{
		Content:  s.getTextMessage(item.MsgType, option.Extra),
		Datetime: item.CreatedAt.Format(time.DateTime),
	})

	return nil
}

func (s *Service) CreateGroupSysMessage(ctx context.Context, option CreateGroupSysMessageOption) error {
	return s.CreateGroupMessage(ctx, CreateGroupMessageOption{
		MsgType:  entity.ChatMsgSysText,
		FromId:   0, // 0:系统消息
		ToFromId: option.GroupId,
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content: option.Content,
		}),
	})
}
