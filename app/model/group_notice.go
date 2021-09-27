package model

type GroupNotice struct {
	ID           int    `json:"id"`
	GroupId      string `json:"group_id"`
	CreatorId    string `json:"creator_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	IsTop        string `json:"is_top"`
	IsDelete     string `json:"is_delete"`
	IsConfirm    string `json:"is_confirm"`
	ConfirmUsers string `json:"confirm_users"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}
