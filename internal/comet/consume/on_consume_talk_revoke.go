package consume

import (
	"context"
	"encoding/json"
	"log/slog"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
)

// 撤销聊天消息
func (h *Handler) onConsumeTalkRevoke(ctx context.Context, body []byte) {
	var in entity.SubEventTalkRevokePayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkRevoke Unmarshal err: %s", err.Error())
		return
	}

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
			data := Message(entity.PushEventImMessageRevoke, entity.ImMessageRevokePayload{
				TalkMode: entity.ChatPrivateMode,
				FromId:   record.FromId,
				ToFromId: record.ToFromId,
				MsgId:    record.MsgId,
				Remark:   in.Remark,
			})

			for _, session := range h.serv.SessionManager().GetSessions(int64(record.UserId)) {
				if err := session.Write(data); err != nil {
					slog.Error("session write message error", "error", err)
				}
			}
		}
	} else if in.TalkMode == entity.ChatGroupMode {
		record, err := h.TalkRecordsService.FindTalkGroupRecord(ctx, in.MsgId)
		if err != nil {
			logger.Errorf("onConsumeTalkRevoke FindTalkGroupRecord err: %s", err.Error())
			return
		}

		data := Message(entity.PushEventImMessageRevoke, entity.ImMessageRevokePayload{
			TalkMode: record.TalkMode,
			FromId:   record.FromId,
			ToFromId: record.ToFromId,
			MsgId:    record.MsgId,
			Remark:   in.Remark,
		})

		for _, uid := range h.GroupMemberRepo.GetMemberIds(ctx, record.ToFromId) {
			for _, session := range h.serv.SessionManager().GetSessions(int64(uid)) {
				if err := session.Write(data); err != nil {
					slog.Error("session write message error", "error", err)
				}
			}
		}
	}
}
