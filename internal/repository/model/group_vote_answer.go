package model

import "time"

type GroupVoteAnswer struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 答题ID
	VoteId    int       `gorm:"column:vote_id;" json:"vote_id"`                 // 投票ID
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	Option    string    `gorm:"column:option;" json:"option"`                   // 投票选项[A、B、C 、D、E、F]
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 答题时间
}

func (GroupVoteAnswer) TableName() string {
	return "group_vote_answer"
}
