package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/model"
	"go-chat/app/pkg/im"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/service"
	"go-chat/config"
	"strconv"
)

type MessagePayload struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type WsSubscribe struct {
	rds                *redis.Client
	conf               *config.Config
	talkRecordsService *service.TalkRecordsService
	ws                 *cache.WsClientSession
	room               *cache.GroupRoom
	contactService     *service.ContactService
}

func NewWsSubscribe(rds *redis.Client, conf *config.Config, talkRecordsService *service.TalkRecordsService, ws *cache.WsClientSession, room *cache.GroupRoom, contactService *service.ContactService) *WsSubscribe {
	return &WsSubscribe{rds: rds, conf: conf, talkRecordsService: talkRecordsService, ws: ws, room: room, contactService: contactService}
}

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

func (w *WsSubscribe) Handle(ctx context.Context) error {
	gateway := fmt.Sprintf(entity.SubscribeWsGatewayPrivate, w.conf.GetSid())

	channels := []string{
		entity.SubscribeWsGatewayAll, // 全局通道
		gateway,                      // 私有通道
		entity.SubscribeCreateGroup,
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	defer sub.Close()

	go func() {
		for msg := range sub.Channel() {
			fmt.Printf("消息订阅 : channel=%s message=%s\n", msg.Channel, msg.Payload)

			switch msg.Channel {
			case entity.SubscribeCreateGroup:
				w.joinGroupRoom(msg.Payload)
			case entity.SubscribeWsGatewayAll, gateway:
				var body *SubscribeBody

				if err := json.Unmarshal([]byte(msg.Payload), &body); err != nil {
					continue
				}

				switch body.EventName {
				case entity.EventTalk:
					w.onConsumeTalk(body.Data)
				case entity.EventKeyboard:
					w.onConsumeKeyboard(body.Data)
				case entity.EventOnlineStatus:
					w.onConsumeOnline(body.Data)
				case entity.EventRevokeTalk:
					w.onConsumeRevokeTalk(body.Data)
				}
			}
		}
	}()

	<-ctx.Done()

	return nil
}

// onConsumeTalk 对话聊天消息
func (w *WsSubscribe) onConsumeTalk(value string) {
	var msg *TalkMessageBody
	if err := json.Unmarshal([]byte(value), &msg); err != nil {
		fmt.Println("onConsumeTalk json", err)
		return
	}

	ctx := context.Background()

	cids := make([]int64, 0)
	if msg.TalkType == 1 {
		arr := [2]int64{msg.SenderID, msg.ReceiverID}
		for _, val := range arr {
			ids := w.ws.GetUidFromClientIds(ctx, w.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(int(val)))

			cids = append(cids, ids...)
		}
	} else {
		ids := w.room.All(ctx, w.conf.GetSid(), strconv.Itoa(int(msg.ReceiverID)))
		cids = append(cids, ids...)
	}

	data, err := w.talkRecordsService.GetTalkRecord(context.Background(), msg.RecordID)
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
func (w *WsSubscribe) onConsumeKeyboard(value string) {
	var msg *KeyboardMessage

	if err := json.Unmarshal([]byte(value), &msg); err != nil {
		return
	}

	cids := w.ws.GetUidFromClientIds(context.Background(), w.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(msg.ReceiverID))

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
func (w *WsSubscribe) onConsumeOnline(value string) {
	var msg *LoginMessage

	if err := json.Unmarshal([]byte(value), &msg); err != nil {
		return
	}

	cids := make([]int64, 0)

	uids := w.contactService.GetContactIds(context.Background(), msg.UserID)
	sid := w.conf.GetSid()
	for _, uid := range uids {
		ids := w.ws.GetUidFromClientIds(context.Background(), sid, im.Sessions.Default.Name, fmt.Sprintf("%d", uid))

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
func (w *WsSubscribe) onConsumeRevokeTalk(value string) {
	var (
		msg    map[string]int
		record *model.TalkRecords
		ctx    = context.Background()
	)

	if err := jsonutil.JsonDecode(value, &msg); err != nil {
		return
	}

	if _, ok := msg["record_id"]; !ok {
		return
	}

	if err := w.talkRecordsService.Db().First(&record, msg["record_id"]).Error; err != nil {
		return
	}

	cids := make([]int64, 0)
	if record.TalkType == entity.PrivateChat {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := w.ws.GetUidFromClientIds(ctx, w.conf.GetSid(), im.Sessions.Default.Name, strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else {
		cids = w.room.All(ctx, w.conf.GetSid(), strconv.Itoa(record.ReceiverId))
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

// onConsumeFriendApply 好友申请消息
func (w *WsSubscribe) onConsumeFriendApply(value string) {

}

// 加入群聊
func (w *WsSubscribe) joinGroupRoom(value string) {
	var (
		ctx = context.Background()
		sid = w.conf.GetSid()
		m   JoinGroup
	)

	if err := json.Unmarshal([]byte(value), &m); err != nil {
		return
	}

	for _, uid := range m.Uids {
		cids := w.ws.GetUidFromClientIds(ctx, sid, im.Sessions.Default.Name, strconv.Itoa(uid))

		for _, cid := range cids {
			_ = w.room.Add(ctx, w.conf.GetSid(), strconv.Itoa(m.GroupID), cid)
		}
	}
}
