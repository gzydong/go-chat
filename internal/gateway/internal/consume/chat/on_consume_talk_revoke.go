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

	var data ConsumeTalkRevoke
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: ", err.Error())
		return
	}

	var record model.TalkRecords
	if err := h.recordsService.Db().First(&record, data.RecordId).Error; err != nil {
		return
	}

	cids := make([]int64, 0)
	if record.TalkType == entity.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(uid))
			cids = append(cids, ids...)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		cids = h.roomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      h.config.ServerId(),
		})
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.IsAck = true
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: "im.message.revoke",
		Content: map[string]any{
			"talk_type":   record.TalkType,
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"record_id":   record.Id,
		},
	})

	socket.Session.Chat.Write(c)
}
