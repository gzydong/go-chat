package model

type TalkRecordsCode struct {
	ID        int    `json:"id"`
	RecordId  string `json:"record_id"`
	UserId    string `json:"user_id"`
	CodeLang  string `json:"code_lang"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
}
