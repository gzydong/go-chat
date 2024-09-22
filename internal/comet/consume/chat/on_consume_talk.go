package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
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
		fmt.Println("[ChatSubscribe] onConsumeTalkPrivateMessage Unmarshal err:", err.Error())
		return
	}

	// 没在线则不推送
	clientIds, _ := h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), message.UserId)
	if len(clientIds) == 0 {
		return
	}

	var extra any
	if err := json.Unmarshal([]byte(message.Extra), &extra); err != nil {
		return
	}

	var quote any
	if err := json.Unmarshal([]byte(message.Quote), &quote); err != nil {
		return
	}

	body := entity.ImMessagePayloadBody{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		UserId:    message.FromId,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: model.No,
		SendTime:  message.CreatedAt.Format(time.DateTime),
		Extra:     extra,
		Quote:     quote,
	}

	if body.UserId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		body.Nickname = user.Nickname
		body.Avatar = user.Avatar
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)

	c.SetMessage(entity.PushEventImMessage, entity.ImMessagePayload{
		TalkMode: entity.ChatPrivateMode,
		ToFromId: message.ToFromId,
		FromId:   message.FromId,
		Body:     body,
	})

	socket.Session.Chat.Write(c)
}

// 群消息
func (h *Handler) onConsumeTalkGroupMessage(ctx context.Context, in entity.SubEventImMessagePayload) {
	message := model.TalkGroupMessage{}
	if err := json.Unmarshal([]byte(in.Message), &message); err != nil {
		return
	}

	clientIds := h.RoomStorage.GetClientIDAll(int32(message.GroupId))

	if len(clientIds) == 0 {
		return
	}

	var extra any
	if err := json.Unmarshal([]byte(message.Extra), &extra); err != nil {
		return
	}

	var quote any
	if err := json.Unmarshal([]byte(message.Quote), &quote); err != nil {
		return
	}

	data := service.TalkRecord{
		MsgId:     message.MsgId,
		Sequence:  int(message.Sequence),
		MsgType:   message.MsgType,
		UserId:    message.FromId,
		Nickname:  "",
		Avatar:    "",
		IsRevoked: model.No,
		SendTime:  message.SendTime.Format(time.DateTime),
		Extra:     extra,
		Quote:     quote,
	}

	if data.UserId > 0 {
		user, err := h.UserRepo.FindByIdWithCache(ctx, message.FromId)
		if err != nil {
			return
		}

		data.Nickname = user.Nickname
		data.Avatar = user.Avatar
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)

	c.SetMessage(entity.PushEventImMessage, entity.ImMessagePayload{
		TalkMode: entity.ChatGroupMode,
		ToFromId: message.GroupId,
		FromId:   message.FromId,
		Body:     data,
	})

	socket.Session.Chat.Write(c)
}
