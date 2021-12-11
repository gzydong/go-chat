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

type onConsumeFunc func(data string)

type SubscribeConsume struct {
	conf           *config.Config
	ws             *cache.WsClientSession
	room           *cache.Room
	recordsService *service.TalkRecordsService
	contactService *service.ContactService
}

func NewSubscribeConsume(conf *config.Config, ws *cache.WsClientSession, room *cache.Room, recordsService *service.TalkRecordsService, contactService *service.ContactService) *SubscribeConsume {
	return &SubscribeConsume{conf: conf, ws: ws, room: room, recordsService: recordsService, contactService: contactService}
}

func (s *SubscribeConsume) Handle(event string, data string) {

	handler := make(map[string]onConsumeFunc)

	// 注册消息回调事件
	handler[entity.EventTalk] = s.onConsumeTalk
	handler[entity.EventKeyboard] = s.onConsumeKeyboard
	handler[entity.EventOnlineStatus] = s.onConsumeOnline
	handler[entity.EventRevokeTalk] = s.onConsumeRevokeTalk
	handler[entity.EventJoinGroupRoom] = s.onConsumeAddGroupRoom

	if f, ok := handler[event]; ok {
		f(data)
	} else {
		fmt.Printf("Event: [%s]未注册回调方法\n", event)
	}
}

func (s *SubscribeConsume) onConsumeTalk(body string) {
	var msg struct {
		TalkType   int   `json:"talk_type"`
		SenderID   int64 `json:"sender_id"`
		ReceiverID int64 `json:"receiver_id"`
		RecordID   int64 `json:"record_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		fmt.Println("onConsumeTalk json", err)
		return
	}

	ctx := context.Background()

	cids := make([]int64, 0)
	if msg.TalkType == 1 {
		arr := [2]int64{msg.SenderID, msg.ReceiverID}
		for _, val := range arr {
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.GetSid(), im.Sessions.Default.Name(), strconv.Itoa(int(val)))

			cids = append(cids, ids...)
		}
	} else {
		ids := s.room.All(ctx, &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(int(msg.ReceiverID)),
			Sid:      s.conf.GetSid(),
		})

		cids = append(cids, ids...)
	}

	data, err := s.recordsService.GetTalkRecord(ctx, msg.RecordID)
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
		Content: gin.H{
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
	var msg struct {
		SenderID   int `json:"sender_id"`
		ReceiverID int `json:"receiver_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		return
	}

	cids := s.ws.GetUidFromClientIds(context.Background(), s.conf.GetSid(), im.Sessions.Default.Name(), strconv.Itoa(msg.ReceiverID))

	if len(cids) == 0 {
		return
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event: entity.EventKeyboard,
		Content: gin.H{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		},
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeOnline 用户上线或下线消息
func (s *SubscribeConsume) onConsumeOnline(body string) {
	var msg struct {
		Status int `json:"status"`
		UserID int `json:"user_id"`
	}

	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		return
	}

	ctx := context.Background()
	cids := make([]int64, 0)

	uids := s.contactService.GetContactIds(ctx, msg.UserID)
	sid := s.conf.GetSid()
	for _, uid := range uids {
		ids := s.ws.GetUidFromClientIds(ctx, sid, im.Sessions.Default.Name(), fmt.Sprintf("%d", uid))

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
		msg struct {
			RecordId int `json:"record_id"`
		}
		record *model.TalkRecords
		ctx    = context.Background()
	)

	if err := jsonutil.JsonDecode(body, &msg); err != nil {
		return
	}

	if err := s.recordsService.Db().First(&record, msg.RecordId).Error; err != nil {
		return
	}

	cids := make([]int64, 0)
	if record.TalkType == entity.PrivateChat {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.GetSid(), im.Sessions.Default.Name(), strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else {
		cids = s.room.All(ctx, &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      s.conf.GetSid(),
		})
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
			"record_id":   record.Id,
		},
	})

	im.Sessions.Default.PushSendChannel(c)
}

// nolint onConsumeContactApply 好友申请消息
func (s *SubscribeConsume) onConsumeContactApply(body string) {

}

// onConsumeAddGroupRoom 加入群房间
func (s *SubscribeConsume) onConsumeAddGroupRoom(body string) {
	var (
		ctx = context.Background()
		sid = s.conf.GetSid()
		m   struct {
			GroupID int   `json:"group_id"`
			Uids    []int `json:"uids"`
		}
	)

	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return
	}

	for _, uid := range m.Uids {
		cids := s.ws.GetUidFromClientIds(ctx, sid, im.Sessions.Default.Name(), strconv.Itoa(uid))

		for _, cid := range cids {
			_ = s.room.Add(ctx, &cache.RoomOption{
				Channel:  im.Sessions.Default.Name(),
				RoomType: entity.RoomGroupChat,
				Number:   strconv.Itoa(m.GroupID),
				Sid:      s.conf.GetSid(),
				Cid:      cid,
			})
		}
	}
}
