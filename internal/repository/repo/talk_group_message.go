package repo

import (
	"context"

	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkGroupMessage struct {
	core.Repo[model.TalkGroupMessage]
}

func NewTalkRecordGroup(db *gorm.DB) *TalkGroupMessage {
	return &TalkGroupMessage{Repo: core.NewRepo[model.TalkGroupMessage](db)}
}

func (t *TalkGroupMessage) FindByMsgId(ctx context.Context, msgId string) (*model.TalkGroupMessage, error) {
	return t.FindByWhere(ctx, "msg_id =?", msgId)
}
