package model

type TalkRecords struct {
	ID         int    `json:"id"`
	TalkType   string `json:"talk_type"`
	MsgType    string `json:"msg_type"`
	UserId     string `json:"user_id"`
	ReceiverId string `json:"receiver_id"`
	IsRevoke   string `json:"is_revoke"`
	IsMark     string `json:"is_mark"`
	IsRead     string `json:"is_read"`
	QuoteId    string `json:"quote_id"`
	WarnUsers  string `json:"warn_users"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
