package model

type TalkRecordsInvite struct {
	ID            int    `json:"id"`
	RecordId      string `json:"record_id"`
	Type          string `json:"type"`
	OperateUserId string `json:"operate_user_id"`
	UserIds       string `json:"user_ids"`
}
