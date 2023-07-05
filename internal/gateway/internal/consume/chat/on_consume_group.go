package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/repository/cache"
)

type ConsumeGroupJoin struct {
	Gid  int   `json:"group_id"`
	Type int   `json:"type"`
	Uids []int `json:"uids"`
}

// 加入群房间
func (h *Handler) onConsumeGroupJoin(ctx context.Context, body []byte) {

	var in ConsumeGroupJoin
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Error("[ChatSubscribe] onConsumeGroupJoin Unmarshal err: ", err.Error())
		return
	}

	sid := h.config.ServerId()
	for _, uid := range in.Uids {
		ids := h.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.Itoa(uid))

		for _, cid := range ids {
			opt := &cache.RoomOption{
				Channel:  socket.Session.Chat.Name(),
				RoomType: entity.RoomImGroup,
				Number:   strconv.Itoa(in.Gid),
				Sid:      h.config.ServerId(),
				Cid:      cid,
			}

			if in.Type == 2 {
				_ = h.roomStorage.Del(ctx, opt)
			} else {
				_ = h.roomStorage.Add(ctx, opt)
			}
		}
	}
}
