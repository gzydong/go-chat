package consume

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/model"
)

// 聊天消息事件
func (h *Handler) onConsumeTalk(ctx context.Context, body []byte) {
	var in entity.SubEventImMessagePayload
	if err := json.Unmarshal(body, &in); err != nil {
		fmt.Println("Err SubEventImMessagePayload===>", err)
		logger.Errorf("[ChatSubscribe] onConsumeTalk Unmarshal err: %s", err.Error())
		return
	}

	if in.TalkMode == entity.ChatPrivateMode {
		h.onConsumeTalkPrivateMessage(ctx, in)
	} else if in.TalkMode == entity.ChatGroupMode {
		h.onConsumeTalkGroupMessage(ctx, in)
	}
}

// 私有消息(点对点消息)
func (h *Handler) onConsumeTalkPrivateMessage(ctx context.Context, in entity.SubEventImMessagePayload) {
	message := model.TalkUserMessage{}
	if err := json.Unmarshal([]byte(in.Message), &message); err != nil {
		return
	}

	sessions := h.serv.SessionManager().GetSessions(int64(message.UserId))
	if len(sessions) == 0 {
		return
	}

	body := entity.ImMessagePayloadBody{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		FromId:    message.FromId,
		IsRevoked: message.IsRevoked,
		SendTime:  message.CreatedAt.Format(time.DateTime),
		Extra:     message.Extra,
		Quote:     message.Quote,
	}

	if body.FromId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		body.Nickname = user.Nickname
		body.Avatar = user.Avatar
	}

	data := Message(entity.PushEventImMessage, entity.ImMessagePayload{
		TalkMode: entity.ChatPrivateMode,
		ToFromId: message.ToFromId,
		FromId:   message.FromId,
		Body:     body,
	})

	for _, session := range sessions {
		if err := session.Write(data); err != nil {
			slog.Error("session write message error", "error", err)
		}
	}
}

// 群消息
func (h *Handler) onConsumeTalkGroupMessage(ctx context.Context, in entity.SubEventImMessagePayload) {
	message := model.TalkGroupMessage{}
	if err := json.Unmarshal([]byte(in.Message), &message); err != nil {
		return
	}

	memberIds := h.GroupMemberRepo.GetMemberIds(ctx, message.GroupId)
	if len(memberIds) == 0 {
		return
	}

	var clientIds []int64
	for _, memberId := range memberIds {
		ids := h.serv.SessionManager().GetConnIds(int64(memberId))
		if len(ids) == 0 {
			continue
		}

		clientIds = append(clientIds, ids...)
	}

	data := entity.ImMessagePayloadBody{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		FromId:    message.FromId,
		IsRevoked: message.IsRevoked,
		SendTime:  message.SendTime.Format(time.DateTime),
		Extra:     message.Extra,
		Quote:     message.Quote,
	}

	if data.FromId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		data.Nickname = user.Nickname
		data.Avatar = user.Avatar
	}

	msg := Message(entity.PushEventImMessage, entity.ImMessagePayload{
		TalkMode: entity.ChatGroupMode,
		ToFromId: message.GroupId,
		FromId:   message.FromId,
		Body:     data,
	})

	for _, cid := range clientIds {
		session, err := h.serv.SessionManager().GetSession(cid)
		if err != nil {
			continue
		}

		if err := session.Write(msg); err != nil {
			slog.Error("session write message error", "error", err)
		}
	}
}
