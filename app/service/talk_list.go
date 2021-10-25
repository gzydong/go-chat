package service

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type TalkListService struct {
	db *gorm.DB
}

func NewTalkListService(db *gorm.DB) *TalkListService {
	return &TalkListService{db: db}
}

func (s *TalkListService) IsDisturb(uid int, receiverId int, talkType int) bool {

	result := &model.TalkList{}

	s.db.Model(&model.TalkList{}).Select("is_disturb").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result)

	return result.IsDisturb == 1
}
