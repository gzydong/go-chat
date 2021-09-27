package model

type UsersEmoticon struct {
	ID          int    `json:"id"`
	UserId      string `json:"user_id"`
	EmoticonIds string `json:"emoticon_ids"`
}
