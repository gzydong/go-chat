package dao

import (
	"go-chat/app/model"
)

type TalkListDao struct {
	*BaseDao
}

func NewTalkListDao(base *BaseDao) *TalkListDao {
	return &TalkListDao{base}
}

func (s *TalkListDao) IsDisturb(uid int, receiverId int, talkType int) bool {

	result := &model.TalkList{}

	s.Db().Model(&model.TalkList{}).Select("is_disturb").Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, receiverId, talkType).First(result)

	return result.IsDisturb == 1
}
