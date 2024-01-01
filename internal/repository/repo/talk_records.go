package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecords struct {
	ichat.Repo[model.TalkRecords]
}

func NewTalkRecords(db *gorm.DB) *TalkRecords {
	return &TalkRecords{Repo: ichat.NewRepo[model.TalkRecords](db)}
}

func (t *TalkRecords) FindByMsgId(ctx context.Context, msgId string) (*model.TalkRecords, error) {
	return t.FindByWhere(ctx, "msg_id =?", msgId)
}
