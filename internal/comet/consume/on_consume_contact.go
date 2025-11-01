package consume

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/logger"
	"github.com/gzydong/go-chat/internal/pkg/sliceutil"
	"github.com/gzydong/go-chat/internal/repository/model"
)

// 用户上线或下线消息
func (h *Handler) onConsumeContactStatus(ctx context.Context, body []byte) {
	var in entity.SubEventContactStatusPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactStatus Unmarshal err: %s", err.Error())
		return
	}

	contactIds := h.ContactService.GetContactIds(ctx, in.UserId)
	if isOk, _ := h.OrganizeRepo.IsQiyeMember(ctx, in.UserId); isOk {
		ids, _ := h.OrganizeRepo.GetMemberIds(ctx)
		contactIds = append(contactIds, ids...)
	}

	data := Message(entity.PushEventContactStatus, in)
	for _, uid := range sliceutil.Unique(contactIds) {
		for _, session := range h.serv.SessionManager().GetSessions(uid) {
			if err := session.Write(data); err != nil {
				slog.Error("session write message error", "error", err)
			}
		}
	}
}

// 好友申请消息
func (h *Handler) onConsumeContactApply(ctx context.Context, body []byte) {
	var in entity.SubEventContactApplyPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactApply Unmarshal err: %s", err.Error())
		return
	}

	var apply model.ContactApply
	if err := h.Source.Db().First(&apply, in.ApplyId).Error; err != nil {
		return
	}

	var user model.Users
	if err := h.Source.Db().First(&user, apply.UserId).Error; err != nil {
		return
	}

	data := Message(entity.PushEventContactApply, entity.ImContactApplyPayload{
		UserId:    user.Id,
		Nickname:  user.Nickname,
		Remark:    apply.Remark,
		ApplyTime: apply.CreatedAt.Format(time.DateTime),
	})

	for _, session := range h.serv.SessionManager().GetSessions(int64(apply.FriendId)) {
		if err := session.Write(data); err != nil {
			slog.Error("session write message error", "error", err)
		}
	}
}
