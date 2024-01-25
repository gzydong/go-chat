package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkUserMessage struct {
	core.Repo[model.TalkUserMessage]
}

func NewTalkRecordFriend(db *gorm.DB) *TalkUserMessage {
	return &TalkUserMessage{Repo: core.NewRepo[model.TalkUserMessage](db)}
}

func (t *TalkUserMessage) FindByMsgId(ctx context.Context, msgId string) (*model.TalkUserMessage, error) {
	return t.FindByWhere(ctx, "msg_id = ?", msgId)
}
