package chat

import (
	"context"
	"log"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
)

type TalkReadMessage struct {
	Event   string `json:"event"`
	Content struct {
		MsgIds     []int `json:"msg_id"`
		ReceiverId int   `json:"receiver_id"`
	} `json:"content"`
}

// onReadMessage 消息已读事件
func (h *Handler) onReadMessage(ctx context.Context, client socket.IClient, data []byte) {

	var in TalkReadMessage
	if err := jsonutil.Decode(data, &in); err != nil {
		log.Println("Chat onReadMessage Err: ", err)
		return
	}

	h.memberService.Db().Model(&model.TalkRecords{}).
		Where("id in ? and receiver_id = ? and is_read = 0", in.Content.MsgIds, client.Uid()).
		Update("is_read", 1)

	h.redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventImMessageRead,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   client.Uid(),
			"receiver_id": in.Content.ReceiverId,
			"ids":         in.Content.MsgIds,
		}),
	}))
}
