package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkRecordFriend struct {
	ichat.Repo[model.TalkRecordFriend]
}

func NewTalkRecordFriend(db *gorm.DB) *TalkRecordFriend {
	return &TalkRecordFriend{Repo: ichat.NewRepo[model.TalkRecordFriend](db)}
}

func (t *TalkRecordFriend) FindByMsgId(ctx context.Context, msgId string) (*model.TalkRecordFriend, error) {
	return t.FindByWhere(ctx, "msg_id =?", msgId)
}
