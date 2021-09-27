package model

type TalkRecordsLogin struct {
	ID        int    `json:"id"`
	RecordId  string `json:"record_id"`
	UserId    string `json:"user_id"`
	Ip        string `json:"ip"`
	Platform  string `json:"platform"`
	Agent     string `json:"agent"`
	Address   string `json:"address"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}
