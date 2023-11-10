package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type TalkSession struct {
	ichat.Repo[model.TalkSession]
}

func NewTalkSession(db *gorm.DB) *TalkSession {
	return &TalkSession{Repo: ichat.NewRepo[model.TalkSession](db)}
}

func (t *TalkSession) IsDisturb(uid int, receiverId int, talkType int) bool {
	resp, err := t.Repo.FindByWhere(context.TODO(), "user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType)
	return err == nil && resp.IsDisturb == 1
}

func (t *TalkSession) FindBySessionId(uid int, receiverId int, talkType int) int {

	resp, err := t.Repo.FindByWhere(context.TODO(), "user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType)
	if err != nil {
		return 0
	}

	return resp.Id
}
