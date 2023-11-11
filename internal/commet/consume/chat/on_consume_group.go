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

type ConsumeGroupJoin struct {
	Gid  int   `json:"group_id"`
	Type int   `json:"type"`
	Uids []int `json:"uids"`
}

type ConsumeGroupApply struct {
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
}

// 加入群房间
func (h *Handler) onConsumeGroupJoin(ctx context.Context, body []byte) {

	var in ConsumeGroupJoin
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupJoin Unmarshal err: %s", err.Error())
		return
	}

	sid := h.Config.ServerId()
	for _, uid := range in.Uids {
		ids := h.ClientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.Itoa(uid))

		for _, cid := range ids {
			opt := &cache.RoomOption{
				Channel:  socket.Session.Chat.Name(),
				RoomType: entity.RoomImGroup,
				Number:   strconv.Itoa(in.Gid),
				Sid:      h.Config.ServerId(),
				Cid:      cid,
			}

			if in.Type == 2 {
				_ = h.RoomStorage.Del(ctx, opt)
			} else {
				_ = h.RoomStorage.Add(ctx, opt)
			}
		}
	}
}

// 入群申请通知
func (h *Handler) onConsumeGroupApply(ctx context.Context, body []byte) {

	var in ConsumeGroupApply
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupApply Unmarshal err: %s", err.Error())
		return
	}

	var groupMember model.GroupMember
	if err := h.Source.Db().First(&groupMember, "group_id = ? and leader = ?", in.GroupId, 2).Error; err != nil {
		return
	}

	var groupDetail model.Group
	if err := h.Source.Db().First(&groupDetail, in.GroupId).Error; err != nil {
		return
	}

	var user model.Users
	if err := h.Source.Db().First(&user, in.UserId).Error; err != nil {
		return
	}

	data := make(map[string]any)
	data["group_name"] = groupDetail.Name
	data["nickname"] = user.Nickname

	clientIds := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(groupMember.UserId))

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventGroupApply, data)

	socket.Session.Chat.Write(c)
}
