package chat

import (
	"context"
	"log"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/model"
)

type TalkReadMessage struct {
	Event   string `json:"event"`
	Content struct {
		MsgIds     []string `json:"msg_ids"`
		ReceiverId int      `json:"receiver_id"`
	} `json:"content"`
}

// onReadMessage 消息已读事件
func (h *Handler) onReadMessage(ctx context.Context, client socket.IClient, data []byte) {
	var in TalkReadMessage
	if err := jsonutil.Decode(data, &in); err != nil {
		log.Println("Chat onReadMessage Err: ", err)
		return
	}

	items := make([]model.TalkRecordsRead, 0, len(in.Content.MsgIds))
	for _, msgId := range in.Content.MsgIds {
		items = append(items, model.TalkRecordsRead{
			MsgId:      msgId,
			UserId:     client.Uid(),
			ReceiverId: in.Content.ReceiverId,
		})
	}

	if err := h.Source.Db().Create(items).Error; err != nil {
		logger.Error("TalkRecordsRead batch creation failed", err.Error())
		return
	}

	h.Redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventImMessageRead,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   client.Uid(),
			"receiver_id": in.Content.ReceiverId,
			"msg_ids":     in.Content.MsgIds,
		}),
	}))
}
