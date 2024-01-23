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
	TalkType   int    `json:"talk_type"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	MsgId      string `json:"msg_id"`
}

// 聊天消息事件
func (h *Handler) onConsumeTalk(ctx context.Context, body []byte) {
	var in ConsumeTalk
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalk Unmarshal err: %s", err.Error())
		return
	}

	if in.TalkType == entity.ChatPrivateMode {
		h.onConsumeTalkPrivateMessage(ctx, in)
	} else if in.TalkType == entity.ChatGroupMode {
		h.onConsumeTalkGroupMessage(ctx, in)
	}
}

// 私有消息(点对点消息)
func (h *Handler) onConsumeTalkPrivateMessage(ctx context.Context, in ConsumeTalk) {
	for _, uid := range [2]int64{in.SenderID, in.ReceiverID} {
		clientIds := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.FormatInt(uid, 10))
		if len(clientIds) == 0 {
			continue
		}

		data, err := h.TalkRecordsService.FindTalkPrivateRecord(ctx, int(uid), in.MsgId)
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
}

// 群消息
func (h *Handler) onConsumeTalkGroupMessage(ctx context.Context, in ConsumeTalk) {
	clientIds := h.RoomStorage.All(ctx, &cache.RoomOption{
		Channel:  socket.Session.Chat.Name(),
		RoomType: entity.RoomImGroup,
		Number:   strconv.Itoa(int(in.ReceiverID)),
		Sid:      h.Config.ServerId(),
	})

	if len(clientIds) == 0 {
		return
	}

	data, err := h.TalkRecordsService.FindTalkGroupRecord(ctx, in.MsgId)
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
