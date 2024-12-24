package chat

import (
	"context"
	"encoding/json"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
)

// 撤销聊天消息
func (h *Handler) onConsumeTalkRevoke(ctx context.Context, body []byte) {
	var in entity.SubEventTalkRevokePayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: %s", err.Error())
		return
	}

	var clientIds []int64
	if in.TalkMode == entity.ChatPrivateMode {
		record, err := h.TalkRecordsService.FindPrivateRecordByMsgId(ctx, in.MsgId)
		if err != nil {
			logger.Errorf("onConsumeTalkRevoke FindPrivateRecordByMsgId err: %s", err.Error())
			return
		}

		if record == nil {
			return
		}

		records, err := h.TalkRecordsService.FindAllPrivateRecordByOriMsgId(ctx, record.OrgMsgId)
		if err != nil {
			logger.Errorf("onConsumeTalkRevoke FindAllPrivateRecordByOriMsgId err: %s", err.Error())
			return
		}

		for _, record := range records {
			clientIds, _ = h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), record.UserId)

			if len(clientIds) == 0 {
				continue
			}

			c := socket.NewSenderContent()
			c.SetAck(true)
			c.SetReceive(clientIds...)
			c.SetMessage(entity.PushEventImMessageRevoke, entity.ImMessageRevokePayload{
				TalkMode: entity.ChatPrivateMode,
				FromId:   record.FromId,
				ToFromId: record.ToFromId,
				MsgId:    record.MsgId,
				Remark:   in.Remark,
			})

			socket.Session.Chat.Write(c)
		}

	} else if in.TalkMode == entity.ChatGroupMode {
		record, err := h.TalkRecordsService.FindTalkGroupRecord(ctx, in.MsgId)
		if err != nil {
			logger.Errorf("onConsumeTalkRevoke FindTalkGroupRecord err: %s", err.Error())
			return
		}

		clientIds = h.RoomStorage.GetClientIDAll(int32(record.ToFromId))
		if len(clientIds) == 0 {
			return
		}

		c := socket.NewSenderContent()
		c.SetAck(true)
		c.SetReceive(clientIds...)
		c.SetMessage(entity.PushEventImMessageRevoke, entity.ImMessageRevokePayload{
			TalkMode: record.TalkMode,
			FromId:   record.FromId,
			ToFromId: record.ToFromId,
			MsgId:    record.MsgId,
			Remark:   in.Remark,
		})

		socket.Session.Chat.Write(c)
	}
}
