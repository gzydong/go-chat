package consume

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/service"
)

type onConsumeFunc func(data string)

type ChatSubscribe struct {
	config         *config.Config
	clientStorage  *cache.ClientStorage
	roomStorage    *cache.RoomStorage
	recordsService *service.TalkRecordsService
	contactService *service.ContactService
	handlers       map[string]onConsumeFunc
}

func NewChatSubscribe(config *config.Config, clientStorage *cache.ClientStorage, roomStorage *cache.RoomStorage, recordsService *service.TalkRecordsService, contactService *service.ContactService) *ChatSubscribe {
	return &ChatSubscribe{config: config, clientStorage: clientStorage, roomStorage: roomStorage, recordsService: recordsService, contactService: contactService}
}

// Events 注册事件
func (s *ChatSubscribe) init() {
	s.handlers = make(map[string]onConsumeFunc)

	s.handlers[entity.EventTalk] = s.onConsumeTalk
	s.handlers[entity.EventTalkKeyboard] = s.onConsumeTalkKeyboard
	s.handlers[entity.EventOnlineStatus] = s.onConsumeLogin
	s.handlers[entity.EventTalkRevoke] = s.onConsumeTalkRevoke
	s.handlers[entity.EventTalkJoinGroup] = s.onConsumeTalkJoinGroup
	s.handlers[entity.EventContactApply] = s.onConsumeContactApply
	s.handlers[entity.EventTalkRead] = s.onConsumeTalkRead
}

// Call 触发回调事件
func (s *ChatSubscribe) Call(event string, data string) {

	if s.handlers == nil {
		s.init()
	}

	if f, ok := s.handlers[event]; ok {
		f(data)
	} else {
		logger.Warnf("ChatSubscribe Event: [%s]未注册回调方法\n", event)
	}
}

// onConsumeTalk 聊天消息事件
func (s *ChatSubscribe) onConsumeTalk(body string) {
	var msg struct {
		TalkType   int   `json:"talk_type"`
		SenderID   int64 `json:"sender_id"`
		ReceiverID int64 `json:"receiver_id"`
		RecordID   int64 `json:"record_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalk Unmarshal err: ", err.Error())
		return
	}

	ctx := context.Background()

	cids := make([]int64, 0)
	if msg.TalkType == entity.ChatPrivateMode {
		for _, val := range [2]int64{msg.SenderID, msg.ReceiverID} {
			ids := s.clientStorage.GetUidFromClientIds(ctx, s.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(int(val)))

			cids = append(cids, ids...)
		}
	} else if msg.TalkType == entity.ChatGroupMode {
		ids := s.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(int(msg.ReceiverID)),
			Sid:      s.config.ServerId(),
		})

		cids = append(cids, ids...)
	}

	data, err := s.recordsService.GetTalkRecord(ctx, msg.RecordID)
	if err != nil {
		logger.Error("[ChatSubscribe] 读取对话记录失败 err: ", err.Error())
		return
	}

	if len(cids) == 0 {
		logger.Error("[ChatSubscribe] cids=0 err: ")
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: entity.EventTalk,
		Content: entity.MapStrAny{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
			"talk_type":   msg.TalkType,
			"data":        data,
		},
	})

	socket.Session.Chat.Write(c)
}

// onConsumeTalkKeyboard 键盘输入事件消息
func (s *ChatSubscribe) onConsumeTalkKeyboard(body string) {
	var msg struct {
		SenderID   int `json:"sender_id"`
		ReceiverID int `json:"receiver_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: ", err.Error())
		return
	}

	cids := s.clientStorage.GetUidFromClientIds(context.Background(), s.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(msg.ReceiverID))

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: entity.EventTalkKeyboard,
		Content: entity.MapStrAny{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		},
	})

	socket.Session.Chat.Write(c)
}

// onConsumeLogin 用户上线或下线消息
func (s *ChatSubscribe) onConsumeLogin(body string) {
	var msg struct {
		Status int `json:"status"`
		UserID int `json:"user_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeLogin Unmarshal err: ", err.Error())
		return
	}

	ctx := context.Background()
	cids := make([]int64, 0)

	uids := s.contactService.GetContactIds(ctx, msg.UserID)
	sid := s.config.ServerId()
	for _, uid := range uids {
		ids := s.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), fmt.Sprintf("%d", uid))

		cids = append(cids, ids...)
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event:   entity.EventOnlineStatus,
		Content: msg,
	})

	socket.Session.Chat.Write(c)
}

// onConsumeTalkRevoke 撤销聊天消息
func (s *ChatSubscribe) onConsumeTalkRevoke(body string) {
	var (
		msg struct {
			RecordId int `json:"record_id"`
		}
		record *model.TalkRecords
		ctx    = context.Background()
	)

	if err := jsonutil.Decode(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: ", err.Error())
		return
	}

	if err := s.recordsService.Db().First(&record, msg.RecordId).Error; err != nil {
		return
	}

	cids := make([]int64, 0)
	if record.TalkType == entity.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := s.clientStorage.GetUidFromClientIds(ctx, s.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		cids = s.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      s.config.ServerId(),
		})
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: entity.EventTalkRevoke,
		Content: entity.MapStrAny{
			"talk_type":   record.TalkType,
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"record_id":   record.Id,
		},
	})

	socket.Session.Chat.Write(c)
}

// nolint onConsumeContactApply 好友申请消息
func (s *ChatSubscribe) onConsumeContactApply(body string) {
	var (
		msg struct {
			ApplId int `json:"apply_id"`
			Type   int `json:"type"`
		}
		ctx = context.Background()
	)

	if err := jsonutil.Decode(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	apply := &model.ContactApply{}
	if err := s.contactService.Db().First(&apply, msg.ApplId).Error; err != nil {
		return
	}

	cids := s.clientStorage.GetUidFromClientIds(ctx, s.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(apply.FriendId))
	if len(cids) == 0 {
		return
	}

	user := &model.Users{}
	if err := s.contactService.Db().First(&user, apply.FriendId).Error; err != nil {
		return
	}

	data := entity.MapStrAny{}
	data["sender_id"] = apply.UserId
	data["receiver_id"] = apply.FriendId
	data["remark"] = apply.Remark
	data["friend"] = entity.MapStrAny{
		"nickname":   user.Nickname,
		"remark":     apply.Remark,
		"created_at": timeutil.FormatDatetime(apply.CreatedAt),
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event:   entity.EventContactApply,
		Content: data,
	})

	socket.Session.Chat.Write(c)
}

// onConsumeTalkJoinGroup 加入群房间
func (s *ChatSubscribe) onConsumeTalkJoinGroup(body string) {
	var (
		ctx  = context.Background()
		sid  = s.config.ServerId()
		data struct {
			Gid  int   `json:"group_id"`
			Type int   `json:"type"`
			Uids []int `json:"uids"`
		}
	)

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkJoinGroup Unmarshal err: ", err.Error())
		return
	}

	for _, uid := range data.Uids {
		cids := s.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.Itoa(uid))

		for _, cid := range cids {
			opts := &cache.RoomOption{
				Channel:  socket.Session.Chat.Name(),
				RoomType: entity.RoomImGroup,
				Number:   strconv.Itoa(data.Gid),
				Sid:      s.config.ServerId(),
				Cid:      cid,
			}

			if data.Type == 2 {
				_ = s.roomStorage.Del(ctx, opts)
			} else {
				_ = s.roomStorage.Add(ctx, opts)
			}
		}
	}
}

// onConsumeTalkRead 消息已读事件
func (s *ChatSubscribe) onConsumeTalkRead(body string) {
	var (
		ctx  = context.Background()
		sid  = s.config.ServerId()
		data struct {
			SenderId   int   `json:"sender_id"`
			ReceiverId int   `json:"receiver_id"`
			Ids        []int `json:"ids"`
		}
	)

	if err := jsonutil.Decode(body, &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	cids := s.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), fmt.Sprintf("%d", data.ReceiverId))

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: entity.EventTalkRead,
		Content: entity.MapStrAny{
			"sender_id":   data.SenderId,
			"receiver_id": data.ReceiverId,
			"ids":         data.Ids,
		},
	})

	socket.Session.Chat.Write(c)
}
