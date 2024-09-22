package message

import (
	"context"
	"github.com/samber/lo"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/model"
	"time"
)

type ForwardMessageOpt struct {
	MsgIds       []string `json:"msg_ids"`
	TalkMode     int      `json:"talk_mode"`
	ToFromId     int      `json:"to_from_id"`
	UserId       int      `json:"user_id"`
	ToUserId     int      `json:"to_user_id"`
	ToUserIdType int      `json:"to_user_id_type"` // 1:用户ID 2:群ID
}

// SplitForward 分拆转发
func (s *Service) toSplitForward(ctx context.Context, req ForwardMessageOpt) error {
	var (
		now          = time.Now()
		db           = s.Source.Db().WithContext(ctx)
		messageItems = make([]model.TalkRecord, 0)
	)

	if req.TalkMode == entity.ChatGroupMode {
		records := make([]model.TalkGroupMessage, 0)

		err := db.Table("talk_group_message").Where("group_id = ? and msg_id in ?", req.ToFromId, req.MsgIds).Scan(&records).Error
		if err != nil {
			return err
		}

		for _, v := range records {
			messageItems = append(messageItems, model.TalkRecord{
				MsgType: v.MsgType,
				Extra:   v.Extra,
			})
		}
	} else {
		records := make([]model.TalkUserMessage, 0)

		err := db.Table("talk_user_message").Where("user_id = ? and to_from_id = ? and msg_id in ?", req.UserId, req.ToFromId, req.MsgIds).Scan(&records).Error
		if err != nil {
			return err
		}

		for _, v := range records {
			messageItems = append(messageItems, model.TalkRecord{
				MsgType: v.MsgType,
				Extra:   v.Extra,
			})
		}
	}

	// 向群发送消息
	if req.ToUserIdType == entity.ChatGroupMode {
		sequences := s.Sequence.BatchGet(ctx, req.ToUserId, false, int64(len(messageItems)))

		items := make([]model.TalkGroupMessage, 0)
		for i, v := range messageItems {
			items = append(items, model.TalkGroupMessage{
				MsgId:     strutil.NewMsgId(),
				Sequence:  sequences[i],
				MsgType:   v.MsgType,
				GroupId:   req.ToUserId,
				FromId:    req.UserId,
				IsRevoked: model.No,
				Extra:     v.Extra,
				Quote:     "{}",
				SendTime:  now,
			})
		}

		if err := db.Create(items).Error; err == nil {
			err = s.PushMessage.MultiPush(ctx, entity.ImTopicChat,
				lo.Map(items, func(item model.TalkGroupMessage, index int) *entity.SubscribeMessage {
					return &entity.SubscribeMessage{
						Event: entity.SubEventImMessage,
						Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
							TalkMode: entity.ChatGroupMode,
							Message:  jsonutil.Encode(item),
						}),
					}
				}),
			)

			if err != nil {
				logger.Errorf("split forward message failed :%s", err.Error())
			}
		}
	} else {
		sequence1 := s.Sequence.BatchGet(ctx, req.ToUserId, true, int64(len(messageItems)))
		sequence2 := s.Sequence.BatchGet(ctx, req.UserId, true, int64(len(messageItems)))

		items := make([]model.TalkUserMessage, 0)
		for i, v := range messageItems {
			msgId := strutil.NewMsgId()

			// 向好友发送消息
			items = append(items, model.TalkUserMessage{
				MsgId:     strutil.NewMsgId(),
				OrgMsgId:  msgId,
				Sequence:  sequence1[i],
				MsgType:   v.MsgType,
				UserId:    req.ToUserId,
				ToFromId:  req.UserId,
				FromId:    req.UserId,
				IsRevoked: model.No,
				IsDeleted: model.No,
				Extra:     v.Extra,
				Quote:     "{}",
				SendTime:  now,
			})

			// 向自己发送消息
			items = append(items, model.TalkUserMessage{
				MsgId:     strutil.NewMsgId(),
				OrgMsgId:  msgId,
				Sequence:  sequence2[i],
				MsgType:   v.MsgType,
				UserId:    req.UserId,
				ToFromId:  req.ToUserId,
				FromId:    req.UserId,
				IsRevoked: model.No,
				IsDeleted: model.No,
				Extra:     v.Extra,
				Quote:     "{}",
				SendTime:  now,
			})
		}

		if err := db.Create(items).Error; err == nil {
			list := lo.Map(items, func(item model.TalkUserMessage, _ int) *entity.SubscribeMessage {
				return &entity.SubscribeMessage{
					Event: entity.SubEventImMessage,
					Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
						TalkMode: entity.ChatPrivateMode,
						Message:  jsonutil.Encode(item),
					}),
				}
			})

			_ = s.PushMessage.MultiPush(ctx, entity.ImTopicChat, list)
		} else {
			logger.Errorf("split forward message failed :%s", err.Error())
		}
	}

	// TODO 待完善需要更新用户未读数

	return nil
}

// CombineForward 合并转发
func (s *Service) toCombineForward(ctx context.Context, req ForwardMessageOpt) error {
	var (
		now = time.Now()
		db  = s.Source.Db().WithContext(ctx)

		extra = model.TalkRecordExtraForward{
			TalkType:   req.TalkMode,
			UserId:     req.UserId,
			ReceiverId: req.ToFromId,
			MsgIds:     req.MsgIds,
			Records:    make([]model.TalkRecordExtraForwardRecord, 0),
		}
		pushMessageItems = make([]entity.SubEventImMessagePayload, 0)
	)

	if req.TalkMode == entity.ChatGroupMode {
		records := make([]model.TalkGroupMessage, 0)

		err := db.Table("talk_group_message").Where("group_id = ? and msg_id in ?", req.ToFromId, req.MsgIds).Order("sequence asc").Limit(3).Scan(&records).Error
		if err != nil {
			return err
		}

		uids := make([]int, 0)
		for _, v := range records {
			uids = append(uids, v.FromId)
		}

		userNameItems, err := s.findUserNameList(ctx, uids)
		if err != nil {
			return err
		}

		for _, v := range records {
			extra.Records = append(extra.Records, model.TalkRecordExtraForwardRecord{
				Nickname: userNameItems[v.FromId],
				Content:  text(v.MsgType, v.Extra),
			})
		}
	} else {
		records := make([]model.TalkUserMessage, 0)

		err := s.Source.Db().Table("talk_user_message").Where("user_id = ? and to_from_id = ? and msg_id in ?", req.UserId, req.ToFromId, req.MsgIds).Order("sequence asc").Limit(3).Scan(&records).Error
		if err != nil {
			return err
		}

		uids := make([]int, 0)
		for _, v := range records {
			uids = append(uids, v.FromId)
		}

		userNameItems, err := s.findUserNameList(ctx, uids)
		if err != nil {
			return err
		}

		for _, v := range records {
			extra.Records = append(extra.Records, model.TalkRecordExtraForwardRecord{
				Nickname: userNameItems[v.FromId],
				Content:  text(v.MsgType, v.Extra),
			})
		}
	}

	switch req.ToUserIdType {
	case entity.ChatPrivateMode: // 好友发送消息
		items := make([]model.TalkUserMessage, 0)

		msgId := strutil.NewMsgId()

		// 向好友发送消息
		items = append(items, model.TalkUserMessage{
			MsgId:     strutil.NewMsgId(),
			OrgMsgId:  msgId,
			Sequence:  s.Sequence.Get(ctx, req.ToUserId, true),
			MsgType:   entity.ChatMsgTypeForward,
			UserId:    req.ToUserId,
			ToFromId:  req.UserId,
			FromId:    req.UserId,
			IsRevoked: model.No,
			IsDeleted: model.No,
			Extra:     jsonutil.Encode(extra),
			Quote:     "{}",
			SendTime:  now,
		})

		// 向自己发送消息
		items = append(items, model.TalkUserMessage{
			MsgId:     strutil.NewMsgId(),
			OrgMsgId:  msgId,
			Sequence:  s.Sequence.Get(ctx, req.UserId, true),
			MsgType:   entity.ChatMsgTypeForward,
			UserId:    req.UserId,
			ToFromId:  req.ToUserId,
			FromId:    req.UserId,
			IsRevoked: model.No,
			IsDeleted: model.No,
			Extra:     jsonutil.Encode(extra),
			Quote:     "{}",
			SendTime:  now,
		})

		if err := db.Create(items).Error; err != nil {
			return err
		}

		for _, item := range items {
			pushMessageItems = append(pushMessageItems, entity.SubEventImMessagePayload{
				TalkMode: entity.ChatPrivateMode,
				Message:  jsonutil.Encode(item),
			})
		}

	case entity.ChatGroupMode: // 向群发送消息
		record := model.TalkGroupMessage{
			MsgId:     strutil.NewMsgId(),
			Sequence:  s.Sequence.Get(ctx, req.ToUserId, false),
			MsgType:   entity.ChatMsgTypeForward,
			GroupId:   req.ToUserId,
			FromId:    req.UserId,
			IsRevoked: model.No,
			Extra:     jsonutil.Encode(extra),
			Quote:     "{}",
			SendTime:  now,
		}

		if err := db.Create(&record).Error; err != nil {
			return err
		}

		pushMessageItems = append(pushMessageItems, entity.SubEventImMessagePayload{
			TalkMode: entity.ChatGroupMode,
			Message:  jsonutil.Encode(record),
		})
	}

	if len(pushMessageItems) > 0 {
		err := s.PushMessage.MultiPush(ctx,
			entity.ImTopicChat,
			lo.Map(pushMessageItems, func(item entity.SubEventImMessagePayload, index int) *entity.SubscribeMessage {
				return &entity.SubscribeMessage{
					Event: entity.SubEventImMessage,
					Payload: jsonutil.Encode(entity.SubEventImMessagePayload{
						TalkMode: item.TalkMode,
						Message:  item.Message,
					}),
				}
			}),
		)

		if err != nil {
			logger.Errorf("forward message failed :%s", err.Error())
		}
	}

	return nil
}

func (s *Service) findUserNameList(ctx context.Context, uids []int) (map[int]string, error) {
	users := make([]model.Users, 0)

	err := s.Source.Db().WithContext(ctx).Find(&users, "id in ?", uids).Error
	if err != nil {
		return nil, err
	}

	items := make(map[int]string)
	for _, v := range users {
		items[v.Id] = v.Nickname
	}

	return items, nil
}

func text(msgType int, extra string) string {
	switch msgType {
	case entity.ChatMsgTypeText:
		data := model.TalkRecordExtraText{}
		if err := jsonutil.Decode(extra, &data); err != nil {
			return ""
		}

		return strutil.MtSubstr(data.Content, 0, 200)
	default:
		if value, ok := entity.ChatMsgTypeMapping[msgType]; ok {
			return value
		}
	}

	return "未知消息"
}
