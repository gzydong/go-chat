package model

type GroupMember struct {
	ID        int    `json:"id"`
	GroupId   string `json:"group_id"`
	UserId    string `json:"user_id"`
	Leader    string `json:"leader"`
	IsMute    string `json:"is_mute"`
	IsQuit    string `json:"is_quit"`
	UserCard  string `json:"user_card"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}
