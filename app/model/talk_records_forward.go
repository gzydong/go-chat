package model

type TalkRecordsForward struct {
	ID        int    `json:"id"`
	RecordId  string `json:"record_id"`
	UserId    string `json:"user_id"`
	RecordsId string `json:"records_id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}
