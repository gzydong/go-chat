package chat

import (
	"context"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/model"
)

type TalkReadMessage struct {
	Event string `json:"event"`
	Data  struct {
		MsgIds     []int `json:"msg_id"`
		ReceiverId int   `json:"receiver_id"`
	} `json:"data"`
}

// OnReadMessage 消息已读事件
func (h *Handler) OnReadMessage(ctx context.Context, client im.IClient, data []byte) {

	var m *TalkReadMessage
	if err := jsonutil.Unmarshal(data, &m); err != nil {
		return
	}

	h.memberService.Db().Model(&model.TalkRecords{}).
		Where("id in ? and receiver_id = ? and is_read = 0", m.Data.MsgIds, client.Uid()).
		Update("is_read", 1)

	h.redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventTalkRead,
		"data": jsonutil.Encode(entity.MapStrAny{
			"sender_id":   client.Uid(),
			"receiver_id": m.Data.ReceiverId,
			"ids":         m.Data.MsgIds,
		}),
	}))
}
