package repo

import (
	"go-chat/internal/repository/model"
)

type TalkSession struct {
	*Base
}

func NewTalkSession(base *Base) *TalkSession {
	return &TalkSession{base}
}

func (t *TalkSession) IsDisturb(uid int, receiverId int, talkType int) bool {

	result := &model.TalkSession{}

	err := t.Db.Model(&model.TalkSession{}).Select("is_disturb").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result).Error
	if err != nil {
		return false
	}

	return result.IsDisturb == 1
}

func (t *TalkSession) FindBySessionId(uid int, receiverId int, talkType int) int {
	result := &model.TalkSession{}

	err := t.Db.Model(&model.TalkSession{}).Select("id").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result).Error
	if err != nil {
		return 0
	}

	return result.Id
}
