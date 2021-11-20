package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/entity"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
	"go-chat/config"
)

type MessagePayload struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type WsSubscribe struct {
	rds                *redis.Client
	conf               *config.Config
	talkRecordsService *service.TalkRecordsService
}

func NewWsSubscribe(rds *redis.Client, conf *config.Config, talkRecordsService *service.TalkRecordsService) *WsSubscribe {
	return &WsSubscribe{rds: rds, conf: conf, talkRecordsService: talkRecordsService}
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
	channels := []string{
		"ws:all",                              // 全局通道
		fmt.Sprintf("ws:%s", w.conf.GetSid()), // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	defer sub.Close()

	go func() {
		for msg := range sub.Channel() {
			var body *SubscribeBody

			if err := json.Unmarshal([]byte(msg.Payload), &body); err != nil {
				continue
			}

			fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)

			switch body.EventName {
			case entity.EventTalk:
				w.onConsumeTalk(body.Data)
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

	uids := make([]int64, 0)
	if msg.TalkType == 1 {
		uids = append(uids, msg.SenderID, msg.ReceiverID)
	} else {
		// 读取群成员ID列表
	}

	data, err := w.talkRecordsService.GetTalkRecord(context.Background(), msg.RecordID)
	if err != nil {
		fmt.Println("GetTalkRecord err", err)
		return
	}

	c := im.NewSenderContent()
	// c.SetReceive(uids...)
	c.SetBroadcast(true)
	c.SetMessage(&im.Message{
		Event: "event_talk",
		Content: map[string]interface{}{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
			"talk_type":   msg.TalkType,
			"data":        data,
		},
	})

	im.Session.DefaultChannel.PushSendChannel(c)
}

// onConsumeKeyboard 键盘输入事件消息
func (w *WsSubscribe) onConsumeKeyboard(value string) {

}

// onConsumeOnline 用户上线或下线消息
func (w *WsSubscribe) onConsumeOnline(value string) {

}

// onConsumeRevokeTalk 撤销聊天消息
func (w *WsSubscribe) onConsumeRevokeTalk(value string) {

}

// onConsumeFriendApply 好友申请消息
func (w *WsSubscribe) onConsumeFriendApply(value string) {

}
