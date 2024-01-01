package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
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

	var record model.TalkRecords
	if err := h.Source.Db().First(&record, "msg_id = ?", in.MsgId).Error; err != nil {
		return
	}

	var clientIds []int64
	if record.TalkType == entity.ChatPrivateMode {
		for _, uid := range [2]int{record.UserId, record.ReceiverId} {
			ids := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(uid))
			clientIds = append(clientIds, ids...)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		clientIds = h.RoomStorage.All(ctx, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(record.ReceiverId),
			Sid:      h.Config.ServerId(),
		})
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
