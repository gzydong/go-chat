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

func (s *Service) CreatePrivateMessage(ctx context.Context, option CreatePrivateMessageOption) error {
	var (
		orgMsgId      = strutil.NewMsgId()
		items         = make([]*model.TalkUserMessage, 0)
		quoteJsonText = "{}"
		now           = time.Now()
	)

	if option.QuoteId != "" {
		quoteRecord := &model.TalkUserMessage{}
		if err := s.Db().First(quoteRecord, "msg_id = ?", option.QuoteId).Error; err != nil {
			return err
		}

		user := &model.Users{}
		if err := s.Db().First(user, "id = ?", quoteRecord.FromId).Error; err != nil {
			return err
		}

		queue := &model.Quote{
			QuoteId: option.QuoteId,
			MsgType: 1,
		}

		queue.Nickname = user.Nickname
		queue.Content = s.getTextMessage(quoteRecord.MsgType, quoteRecord.Extra)
		quoteJsonText = jsonutil.Encode(queue)
	}

	items = append(items, &model.TalkUserMessage{
		MsgId:     strutil.NewMsgId(),
		Sequence:  s.Sequence.Get(ctx, option.FromId, true),
		MsgType:   option.MsgType,
		UserId:    option.FromId,
		ToFromId:  option.ToFromId,
		FromId:    option.FromId,
		Extra:     option.Extra,
		Quote:     quoteJsonText,
		OrgMsgId:  orgMsgId,
		SendTime:  now,
		IsRevoked: model.No,
		IsDeleted: model.No,
	})

	items = append(items, &model.TalkUserMessage{
		MsgId:     strutil.NewMsgId(),
		Sequence:  s.Sequence.Get(ctx, option.ToFromId, true),
		MsgType:   option.MsgType,
		UserId:    option.ToFromId,
		ToFromId:  option.FromId,
		FromId:    option.FromId,
		Extra:     option.Extra,
		Quote:     quoteJsonText,
		OrgMsgId:  orgMsgId,
		SendTime:  now,
		IsRevoked: model.No,
		IsDeleted: model.No,
	})

	if err := s.Db().WithContext(ctx).Create(items).Error; err != nil {
		return err
	}

	// 推送消息
	pipe := s.Source.Redis().Pipeline()
	for _, item := range items {
		content := &entity.SubscribeMessage{
			Event: entity.SubEventImMessage,
			Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
				TalkMode: entity.ChatPrivateMode,
				Message:  jsonutil.Encode(item),
			}),
		}

		pipe.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(content))

		if item.UserId != option.FromId {
			s.UnreadStorage.PipeIncr(ctx, pipe, item.UserId, entity.ChatPrivateMode, item.ToFromId)
		}

		// 更新最后一条消息
		_ = s.MessageStorage.Set(ctx, entity.ChatPrivateMode, item.UserId, item.ToFromId, &cache.LastCacheMessage{
			Content:  s.getTextMessage(item.MsgType, option.Extra),
			Datetime: item.CreatedAt.Format(time.DateTime),
		})
	}

	_, _ = pipe.Exec(ctx)

	return nil
}

func (s *Service) CreateToUserPrivateMessage(ctx context.Context, data *model.TalkUserMessage) error {
	if data.MsgId == "" {
		data.MsgId = strutil.NewMsgId()
	}

	if data.OrgMsgId == "" {
		data.OrgMsgId = data.MsgId
	}

	if data.Sequence <= 0 {
		data.Sequence = s.Sequence.Get(ctx, data.UserId, true)
	}

	if data.Quote == "" {
		data.Quote = "{}"
	}

	if data.SendTime.IsZero() {
		data.SendTime = time.Now()
	}

	data.IsRevoked = model.No
	data.IsDeleted = model.No

	if err := s.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	err := s.PushMessage.Push(ctx, entity.ImTopicChat, &entity.SubscribeMessage{
		Event: entity.SubEventImMessage,
		Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
			TalkMode: entity.ChatPrivateMode,
			Message:  jsonutil.Encode(data),
		}),
	})
	if err != nil {
		logger.Errorf("SendToUserPrivateLetter redis push err:%s", err.Error())
	}

	s.UnreadStorage.Incr(ctx, data.UserId, entity.ChatPrivateMode, data.ToFromId)

	// 更新最后一条消息
	_ = s.MessageStorage.Set(ctx, entity.ChatPrivateMode, data.UserId, data.ToFromId, &cache.LastCacheMessage{
		Content:  s.getTextMessage(data.MsgType, data.Extra),
		Datetime: data.CreatedAt.Format(time.DateTime),
	})

	return nil
}

func (s *Service) CreatePrivateSysMessage(ctx context.Context, option CreatePrivateSysMessageOption) error {
	return s.CreateToUserPrivateMessage(ctx, &model.TalkUserMessage{
		MsgId:    strutil.NewMsgId(),
		Sequence: s.Sequence.Get(ctx, option.FromId, true),
		MsgType:  entity.ChatMsgSysText,
		UserId:   option.FromId,
		ToFromId: option.ToFromId,
		FromId:   0,
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content: option.Content,
		}),
		Quote:    "{}",
		SendTime: time.Now(),
	})
}
