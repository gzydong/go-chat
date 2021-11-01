package dao

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type TalkListDao struct {
	db *gorm.DB
}

func NewTalkListDao(db *gorm.DB) *TalkListDao {
	return &TalkListDao{db}
}

func (s *TalkListDao) IsDisturb(uid int, receiverId int, talkType int) bool {

	result := &model.TalkList{}

	s.db.Model(&model.TalkList{}).Select("is_disturb").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result)

	return result.IsDisturb == 1
}
