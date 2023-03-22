package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/sliceutil"
)

type ConsumeLogin struct {
	Status int `json:"status"`
	UserID int `json:"user_id"`
}

// 用户上线或下线消息
func (h *Handler) onConsumeLogin(ctx context.Context, body []byte) {

	var msg *ConsumeLogin
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeLogin Unmarshal err: ", err.Error())
		return
	}

	cids := make([]int64, 0)

	uids := h.contactService.GetContactIds(ctx, msg.UserID)

	isQiye, _ := h.organize.IsQiyeMember(ctx, msg.UserID)
	if isQiye {
		mids, _ := h.organize.GetMemberIds(ctx)
		uids = append(uids, mids...)
	}

	sid := h.config.ServerId()
	for _, uid := range sliceutil.Unique(uids) {
		ids := h.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.FormatInt(uid, 10))
		cids = append(cids, ids...)
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event:   "im.contact.status",
		Content: msg,
	})

	socket.Session.Chat.Write(c)
}
