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

// 聊天消息事件
func (h *Handler) onConsumeTalk(ctx context.Context, body []byte) {

	var in ConsumeTalk
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalk Unmarshal err: %s", err.Error())
		return
	}

	var clientIds []int64
	if in.TalkType == entity.ChatPrivateMode {
		for _, val := range [2]int64{in.SenderID, in.ReceiverID} {
			ids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.FormatInt(val, 10))

			clientIds = append(clientIds, ids...)
		}
	} else if in.TalkType == entity.ChatGroupMode {
		ids := h.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(int(in.ReceiverID)),
			Sid:      h.config.ServerId(),
		})

		clientIds = append(clientIds, ids...)
	}

	if len(clientIds) == 0 {
		return
	}

	data, err := h.recordsService.GetTalkRecord(ctx, in.RecordID)
	if err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)
	c.SetMessage(entity.PushEventImMessage, map[string]any{
		"sender_id":   in.SenderID,
		"receiver_id": in.ReceiverID,
		"talk_type":   in.TalkType,
		"data":        data,
	})

	socket.Session.Chat.Write(c)
}
