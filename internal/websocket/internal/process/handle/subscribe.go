package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/service"
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
	handler[entity.EventFriendApply] = s.onConsumeContactApply

	if f, ok := handler[event]; ok {
		f(data)
	} else {
		fmt.Printf("Event: [%s]未注册回调方法\n", event)
	}
}

// onConsumeTalk 聊天消息事件
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
		for _, val := range [2]int64{msg.SenderID, msg.ReceiverID} {
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.ServerId(), im.Sessions.Default.Name(), strconv.Itoa(int(val)))

			cids = append(cids, ids...)
		}
	} else {
		ids := s.room.All(ctx, &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(int(msg.ReceiverID)),
			Sid:      s.conf.ServerId(),
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
		Event: entity.EventTalk,
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

	cids := s.ws.GetUidFromClientIds(context.Background(), s.conf.ServerId(), im.Sessions.Default.Name(), strconv.Itoa(msg.ReceiverID))

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
	sid := s.conf.ServerId()
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
			ids := s.ws.GetUidFromClientIds(ctx, s.conf.ServerId(), im.Sessions.Default.Name(), strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else {
		cids = s.room.All(ctx, &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      s.conf.ServerId(),
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
	var (
		msg struct {
			ApplId int `json:"apply_id"`
			Type   int `json:"type"`
		}
		ctx   = context.Background()
		apply *model.ContactApply
	)

	if err := jsonutil.JsonDecode(body, &msg); err != nil {
		return
	}

	err := s.contactService.Db().Model(model.ContactApply{}).First(&apply, msg.ApplId).Error
	if err != nil {
		return
	}

	cids := make([]int64, 0)

	if msg.Type == 1 {
		cids = s.ws.GetUidFromClientIds(ctx, s.conf.ServerId(), im.Sessions.Default.Name(), strconv.Itoa(apply.FriendId))
	} else {
		cids = s.ws.GetUidFromClientIds(ctx, s.conf.ServerId(), im.Sessions.Default.Name(), strconv.Itoa(apply.UserId))
	}

	if len(cids) == 0 {
		return
	}

	data := gin.H{}
	if msg.Type == 1 {
		data["sender_id"] = apply.UserId
		data["receiver_id"] = apply.FriendId
		data["remark"] = apply.Remark
	} else {
		data["sender_id"] = apply.FriendId
		data["receiver_id"] = apply.UserId
		data["remark"] = apply.Remark
		data["status"] = 1
	}

	data["friend"] = gin.H{
		"user_id":  1,
		"avatar":   "$friendInfo->avatar",
		"nickname": "$friendInfo->nickname",
		"mobile":   "$friendInfo->mobile",
	}

	c := im.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&im.Message{
		Event:   entity.EventFriendApply,
		Content: data,
	})

	im.Sessions.Default.PushSendChannel(c)
}

// onConsumeAddGroupRoom 加入群房间
func (s *SubscribeConsume) onConsumeAddGroupRoom(body string) {
	var (
		ctx  = context.Background()
		sid  = s.conf.ServerId()
		data struct {
			Gid  int   `json:"group_id"`
			Type int   `json:"type"`
			Uids []int `json:"uids"`
		}
	)

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		fmt.Println("onConsumeAddGroupRoom Unmarshal err: ", err.Error())
		return
	}

	for _, uid := range data.Uids {
		cids := s.ws.GetUidFromClientIds(ctx, sid, im.Sessions.Default.Name(), strconv.Itoa(uid))

		for _, cid := range cids {
			opts := &cache.RoomOption{
				Channel:  im.Sessions.Default.Name(),
				RoomType: entity.RoomGroupChat,
				Number:   strconv.Itoa(data.Gid),
				Sid:      s.conf.ServerId(),
				Cid:      cid,
			}

			if data.Type == 2 {
				_ = s.room.Del(ctx, opts)
			} else {
				_ = s.room.Add(ctx, opts)
			}
		}
	}
}
