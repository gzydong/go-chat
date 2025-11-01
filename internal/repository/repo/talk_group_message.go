package repo

import (
	"context"

	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
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
