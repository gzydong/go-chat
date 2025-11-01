package consume

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/logger"
	"github.com/gzydong/go-chat/internal/repository/model"
)

// 加入群房间
func (h *Handler) onConsumeGroupJoin(ctx context.Context, body []byte) {
	var in entity.SubEventGroupJoinPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeGroupJoin Unmarshal err: %s", err.Error())
		return
	}

	fmt.Println(in)
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
		ids := h.serv.SessionManager().GetConnIds(int64(member.UserId))
		if len(ids) == 0 {
			continue
		}

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

	msg := Message(entity.PushEventGroupApply, entity.ImGroupApplyPayload{
		GroupId:   groupDetail.Id,
		GroupName: groupDetail.Name,
		UserId:    user.Id,
		Nickname:  user.Nickname,
		Remark:    groupApply.Remark,
		ApplyTime: groupApply.CreatedAt.Format(time.DateTime),
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
