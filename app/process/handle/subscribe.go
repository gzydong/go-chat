package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/model"
	"go-chat/app/pkg/im"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/service"
	"go-chat/config"
	"strconv"
)

type JoinGroup struct {
	GroupID int   `json:"group_id"`
	Uids    []int `json:"uids"`
}

type KeyboardMessage struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

type LoginMessage struct {
	Status int `json:"status"`
	UserID int `json:"user_id"`
}

type SubscribeBody struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type TalkMessageBody struct {
	TalkType   int   `json:"talk_type"`
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	RecordID   int64 `json:"record_id"`
}

type SubscribeConsume struct {
	conf               *config.Config
	talkRecordsService *service.TalkRecordsService
	ws                 *cache.WsClientSession
	room               *cache.GroupRoom
	contactService     *service.ContactService
}

func NewSubscribeConsume(conf *config.Config, talkRecordsService *service.TalkRecordsService, ws *cache.WsClientSession, room *cache.GroupRoom, contactService *service.ContactService) *SubscribeConsume {
	return &SubscribeConsume{conf: conf, talkRecordsService: talkRecordsService, ws: ws, room: room, contactService: contactService}
}

func (s *SubscribeConsume) Handle(event string, data string) {
	switch event {
	case entity.EventTalk:
		s.onConsumeTalk(data)
	case entity.EventKeyboard:
		s.onConsumeKeyboard(data)
	case entity.EventOnlineStatus:
		s.onConsumeOnline(data)
	case entity.EventRevokeTalk:
		s.onConsumeRevokeTalk(data)
	case entity.EventJoinGroupRoom:
		s.onConsumeAddGroupRoom(data)
	}
}

func (s *SubscribeConsume) onConsumeTalk(body string) {
	var msg *TalkMessageBody
	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		fmt.Println("onConsumeTalk json", err)
		return
	}

	ctx := context.Background()

	cids := make([]int64, 0)
	if msg.TalkType == 1 {
		arr := [2]int64{msg.SenderID, msg.ReceiverID}
		for _, val := range arr {
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(int(val)))

			cids = append(cids, ids...)
		}
	} else {
		ids := s.room.All(ctx, s.conf.GetSid(), strconv.Itoa(int(msg.ReceiverID)))
		cids = append(cids, ids...)
	}

	data, err := s.talkRecordsService.GetTalkRecord(context.Background(), msg.RecordID)
	if err != nil {
		fmt.Println("GetTalkRecord err", err)
		return
	}

	if len(cids) == 0 {
		return
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event: "event_talk",
		Content: map[string]interface{}{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
			"talk_type":   msg.TalkType,
			"data":        data,
		},
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeKeyboard 键盘输入事件消息
func (s *SubscribeConsume) onConsumeKeyboard(body string) {
	var msg *KeyboardMessage

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		return
	}

	cids := s.ws.GetUidFromClientIds(context.Background(), s.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(msg.ReceiverID))

	if len(cids) == 0 {
		return
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event: entity.EventKeyboard,
		Content: map[string]interface{}{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		},
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeOnline 用户上线或下线消息
func (s *SubscribeConsume) onConsumeOnline(body string) {
	var msg *LoginMessage

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		return
	}

	cids := make([]int64, 0)

	uids := s.contactService.GetContactIds(context.Background(), msg.UserID)
	sid := s.conf.GetSid()
	for _, uid := range uids {
		ids := s.ws.GetUidFromClientIds(context.Background(), sid, im.Sessions.Default.Name, fmt.Sprintf("%d", uid))

		cids = append(cids, ids...)
	}

	if len(cids) == 0 {
		return
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event:   entity.EventOnlineStatus,
		Content: msg,
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeRevokeTalk 撤销聊天消息
func (s *SubscribeConsume) onConsumeRevokeTalk(body string) {
	var (
		msg    map[string]int
		record *model.TalkRecords
		ctx    = context.Background()
	)

	if err := jsonutil.JsonDecode(body, &msg); err != nil {
		return
	}

	if _, ok := msg["record_id"]; !ok {
		return
	}

	if err := s.talkRecordsService.Db().First(&record, msg["record_id"]).Error; err != nil {
		return
	}

	cids := make([]int64, 0)
	if record.TalkType == entity.PrivateChat {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else {
		cids = s.room.All(ctx, s.conf.GetSid(), strconv.Itoa(record.ReceiverId))
	}

	if len(cids) == 0 {
		return
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event: entity.EventRevokeTalk,
		Content: gin.H{
			"talk_type":   record.TalkType,
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"record_id":   record.ID,
		},
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeContactApply 好友申请消息
func (s *SubscribeConsume) onConsumeContactApply(body string) {

}

// onConsumeAddGroupRoom 加入群房间
func (s *SubscribeConsume) onConsumeAddGroupRoom(body string) {
	var (
		ctx = context.Background()
		sid = s.conf.GetSid()
		m   JoinGroup
	)

	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return
	}

	for _, uid := range m.Uids {
		cids := s.ws.GetUidFromClientIds(ctx, sid, im.Sessions.Default.Name, strconv.Itoa(uid))

		for _, cid := range cids {
			_ = s.room.Add(ctx, s.conf.GetSid(), strconv.Itoa(m.GroupID), cid)
		}
	}
}
