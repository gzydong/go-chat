package model

import "time"

type TalkRecordsVote struct {
	ID           int       `json:"id" grom:"comment:投票ID"`
	RecordId     int       `json:"record_id" grom:"comment:消息记录ID"`
	UserId       int       `json:"user_id" grom:"comment:用户ID"`
	Title        string    `json:"title" grom:"comment:投票标题"`
	AnswerMode   int       `json:"answer_mode" grom:"comment:投票模式"`
	AnswerOption string    `json:"answer_option" grom:"comment:投票选项"`
	AnswerNum    int       `json:"answer_num" grom:"comment:应答人数"`
	AnsweredNum  int       `json:"answered_num" grom:"comment:已答人数"`
	Status       int       `json:"status" grom:"comment:投票状态"`
	CreatedAt    time.Time `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" grom:"comment:更新时间"`
}
