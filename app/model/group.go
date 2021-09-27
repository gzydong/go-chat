package model

type Group struct {
	ID          int    `json:"id"`
	CreatorId   string `json:"creator_id"`
	GroupName   string `json:"group_name"`
	Profile     string `json:"profile"`
	Avatar      string `json:"avatar"`
	MaxNum      string `json:"max_num"`
	IsOvert     string `json:"is_overt"`
	IsMute      string `json:"is_mute"`
	IsDismiss   string `json:"is_dismiss"`
	CreatedAt   string `json:"created_at"`
	DismissedAt string `json:"dismissed_at"`
}
