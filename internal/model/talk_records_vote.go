package model

import "time"

type TalkRecordsVote struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`             // 投票ID
	RecordId     int       `gorm:"column:record_id;default:0;NOT NULL" json:"record_id"`       // 消息记录ID
	UserId       int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`           // 用户ID
	Title        string    `gorm:"column:title;NOT NULL" json:"title"`                         // 投票标题
	AnswerMode   int       `gorm:"column:answer_mode;default:0;NOT NULL" json:"answer_mode"`   // 答题模式[0:单选;1:多选;]
	AnswerOption string    `gorm:"column:answer_option;NOT NULL" json:"answer_option"`         // 答题选项
	AnswerNum    int       `gorm:"column:answer_num;default:0;NOT NULL" json:"answer_num"`     // 应答人数
	AnsweredNum  int       `gorm:"column:answered_num;default:0;NOT NULL" json:"answered_num"` // 已答人数
	IsAnonymous  int       `gorm:"column:is_anonymous;default:0;NOT NULL" json:"is_anonymous"` // 匿名投票[0:否;1:是;]
	Status       int       `gorm:"column:status;default:0;NOT NULL" json:"status"`             // 投票状态[0:投票中;1:已完成;]
	CreatedAt    time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`               // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`               // 更新时间
}

type QueryVoteModel struct {
	RecordId     int    `json:"record_id" grom:"comment:消息记录ID"`
	ReceiverId   int    `json:"receiver_id"`
	TalkType     int    `json:"talk_type"`
	MsgType      int    `json:"msg_type"`
	VoteId       int    `json:"vote_id"`
	AnswerMode   int    `json:"answer_mode"`
	AnswerOption string `json:"answer_option"`
	AnswerNum    int    `json:"answer_num"`
	VoteStatus   int    `json:"vote_status"`
}
