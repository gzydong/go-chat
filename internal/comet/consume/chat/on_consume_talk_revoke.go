package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
	"go-chat/internal/repository/model"
)

type ConsumeTalkRevoke struct {
	MsgId string `json:"msg_id"`
}

// 撤销聊天消息
func (h *Handler) onConsumeTalkRevoke(ctx context.Context, body []byte) {
	var in ConsumeTalkRevoke
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: %s", err.Error())
		return
	}

	var record model.TalkRecord
	if err := h.Source.Db().First(&record, "msg_id = ?", in.MsgId).Error; err != nil {
		return
	}

	var clientIds []int64
	if record.TalkType == entity.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), uid)
			clientIds = append(clientIds, ids...)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		clientIds = h.RoomStorage.GetClientIDAll(int32(record.ReceiverId))
	}

	if len(clientIds) == 0 {
		return
	}

	var user model.Users
	if err := h.Source.Db().WithContext(ctx).Select("id,nickname").First(&user, record.UserId).Error; err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventImMessageRevoke, map[string]any{
		"talk_type":   record.TalkType,
		"sender_id":   record.UserId,
		"receiver_id": record.ReceiverId,
		"msg_id":      record.MsgId,
		"text":        fmt.Sprintf("%s: 撤回了一条消息", user.Nickname),
	})

	socket.Session.Chat.Write(c)
}
