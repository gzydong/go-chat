package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

type ConsumeTalk struct {
	TalkType   int   `json:"talk_type"`
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	RecordID   int64 `json:"record_id"`
}

// onConsumeTalk 聊天消息事件
func (h *Handler) onConsumeTalk(ctx context.Context, body []byte) {

	var msg ConsumeTalk
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalk Unmarshal err: ", err.Error())
		return
	}

	cids := make([]int64, 0)
	if msg.TalkType == entity.ChatPrivateMode {
		for _, val := range [2]int64{msg.SenderID, msg.ReceiverID} {
			ids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.FormatInt(val, 10))

			cids = append(cids, ids...)
		}
	} else if msg.TalkType == entity.ChatGroupMode {
		ids := h.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(int(msg.ReceiverID)),
			Sid:      h.config.ServerId(),
		})

		cids = append(cids, ids...)
	}

	data, err := h.recordsService.GetTalkRecord(ctx, msg.RecordID)
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
