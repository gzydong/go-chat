package model

import "time"

const (
	VoteAnswerModeSingle   = 1 // 单选
	VoteAnswerModeMultiple = 2 // 多选

	VoteStatusWait   = 1 // 等待中
	VoteStatusFinish = 2 // 已结束
)

type GroupVote struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 投票ID
	GroupId      int       `gorm:"column:group_id;" json:"group_id"`               // 群组ID
	UserId       int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	Title        string    `gorm:"column:title;" json:"title"`                     // 投票标题
	AnswerMode   int       `gorm:"column:answer_mode;" json:"answer_mode"`         // 答题模式[1:单选;1:多选;]
	AnswerOption string    `gorm:"column:answer_option;" json:"answer_option"`     // 答题选项
	AnswerNum    int       `gorm:"column:answer_num;" json:"answer_num"`           // 应答人数
	AnsweredNum  int       `gorm:"column:answered_num;" json:"answered_num"`       // 已答人数
	IsAnonymous  int       `gorm:"column:is_anonymous;" json:"is_anonymous"`       // 匿名投票[1:否;2:是;]
	Status       int       `gorm:"column:status;" json:"status"`                   // 投票状态[1:投票中;2:已完成;]
	CreatedAt    time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (GroupVote) TableName() string {
	return "group_vote"
}

type GroupVoteOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueryVoteModel struct {
	MsgId        string `json:"msg_id"`
	ReceiverId   int    `json:"receiver_id"`
	TalkType     int    `json:"talk_type"`
	MsgType      int    `json:"msg_type"`
	VoteId       int    `json:"vote_id"`
	AnswerMode   int    `json:"answer_mode"`
	AnswerOption string `json:"answer_option"`
	AnswerNum    int    `json:"answer_num"`
	VoteStatus   int    `json:"vote_status"`
}
