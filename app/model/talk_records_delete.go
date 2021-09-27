package model

type TalkRecordsDelete struct {
	ID        int    `json:"id"`
	RecordId  string `json:"record_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
