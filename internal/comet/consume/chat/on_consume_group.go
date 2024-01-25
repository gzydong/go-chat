package chat

import (
	"context"
	"encoding/json"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
	"go-chat/internal/repository/model"
)

// 加入群房间
func (h *Handler) onConsumeGroupJoin(ctx context.Context, body []byte) {
	var in entity.SubEventGroupJoinPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupJoin Unmarshal err: %s", err.Error())
		return
	}

	sid := server.ID()
	now := time.Now()
	for _, uid := range in.Uids {
		ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), uid)

		for _, cid := range ids {
			if in.Type == 2 {
				_ = h.RoomStorage.Delete(int32(in.GroupId), cid, now.Unix())
			} else {
				_ = h.RoomStorage.Insert(int32(in.GroupId), cid, now.Unix())
			}
		}
	}
}

// 入群申请通知
func (h *Handler) onConsumeGroupApply(ctx context.Context, body []byte) {
	var in entity.SubEventGroupApplyPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupApply Unmarshal err: %s", err.Error())
		return
	}

	var members []model.GroupMember
	if err := h.Source.Db().Find(&members, "group_id = ? and leader in ? and is_quit = ?", in.GroupId, []int{model.GroupMemberLeaderOwner, model.GroupMemberLeaderAdmin}, model.No).Error; err != nil {
		return
	}

	var clientIds []int64
	for _, member := range members {
		ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), member.UserId)
		clientIds = append(clientIds, ids...)
	}

	if len(clientIds) == 0 {
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

	var groupApply model.GroupApply
	if err := h.Source.Db().First(&groupApply, in.ApplyId).Error; err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetAck(true)
	c.SetMessage(entity.PushEventGroupApply, entity.ImGroupApplyPayload{
		GroupId:   groupDetail.Id,
		GroupName: groupDetail.Name,
		UserId:    user.Id,
		Nickname:  user.Nickname,
		Remark:    groupApply.Remark,
		ApplyTime: groupApply.CreatedAt.Format(time.DateTime),
	})

	socket.Session.Chat.Write(c)
}
