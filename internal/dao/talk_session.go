package dao

import (
	"go-chat/internal/model"
)

type TalkSessionDao struct {
	*BaseDao
}

func NewTalkSessionDao(base *BaseDao) *TalkSessionDao {
	return &TalkSessionDao{base}
}

func (s *TalkSessionDao) IsDisturb(uid int, receiverId int, talkType int) bool {

	result := &model.TalkSession{}

	err := s.Db().Model(&model.TalkSession{}).Select("is_disturb").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result).Error
	if err != nil {
		return false
	}

	return result.IsDisturb == 1
}

func (s *TalkSessionDao) FindBySessionId(uid int, receiverId int, talkType int) int {
	result := &model.TalkSession{}

	err := s.Db().Model(&model.TalkSession{}).Select("id").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result).Error
	if err != nil {
		return 0
	}

	return result.Id
}
