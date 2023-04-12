package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
)

type ConsumeTalkRevoke struct {
	RecordId int `json:"record_id"`
}

// 撤销聊天消息
func (h *Handler) onConsumeTalkRevoke(ctx context.Context, body []byte) {

	var in ConsumeTalkRevoke
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: ", err.Error())
		return
	}

	var record model.TalkRecords
	if err := h.recordsService.Db().First(&record, in.RecordId).Error; err != nil {
		return
	}

	var clientIds []int64
	if record.TalkType == entity.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(uid))
			clientIds = append(clientIds, ids...)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		clientIds = h.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      h.config.ServerId(),
		})
	}

	if len(clientIds) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventImMessageRevoke, map[string]any{
		"talk_type":   record.TalkType,
		"sender_id":   record.UserId,
		"receiver_id": record.ReceiverId,
		"record_id":   record.Id,
	})

	socket.Session.Chat.Write(c)
}
