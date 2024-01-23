package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecordGroup struct {
	ichat.Repo[model.TalkRecordGroup]
}

func NewTalkRecordGroup(db *gorm.DB) *TalkRecordGroup {
	return &TalkRecordGroup{Repo: ichat.NewRepo[model.TalkRecordGroup](db)}
}

func (t *TalkRecordGroup) FindByMsgId(ctx context.Context, msgId string) (*model.TalkRecordGroup, error) {
	return t.FindByWhere(ctx, "msg_id =?", msgId)
}
