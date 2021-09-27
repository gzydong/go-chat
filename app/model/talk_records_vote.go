package model

type TalkRecordsVote struct {
	ID           int    `json:"id"`
	RecordId     string `json:"record_id"`
	UserId       string `json:"user_id"`
	Title        string `json:"title"`
	AnswerMode   string `json:"answer_mode"`
	AnswerOption string `json:"answer_option"`
	AnswerNum    string `json:"answer_num"`
	AnsweredNum  string `json:"answered_num"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
