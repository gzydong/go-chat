package model

type UsersEmoticon struct {
	ID          int    `json:"id" grom:"comment:收藏ID"`
	UserId      int    `json:"user_id" grom:"comment:用户ID"`
	EmoticonIds string `json:"emoticon_ids" grom:"comment:表情包ID，多个用英文逗号拼接"`
}
