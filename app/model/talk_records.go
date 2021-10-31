package model

import "time"

type TalkRecords struct {
	ID         int       `json:"id" grom:"comment:聊天消息ID"`
	TalkType   int       `json:"talk_type" grom:"comment:对话类型"`
	MsgType    int       `json:"msg_type" grom:"comment:消息类型"`
	UserId     int       `json:"user_id" grom:"comment:发送者ID"`
	ReceiverId int       `json:"receiver_id" grom:"comment:接收者ID"`
	IsRevoke   int       `json:"is_revoke" grom:"comment:是否撤回消息"`
	IsMark     int       `json:"is_mark" grom:"comment:是否重要消息"`
	IsRead     int       `json:"is_read" grom:"comment:是否已读"`
	QuoteId    int       `json:"quote_id" grom:"comment:引用消息ID"`
	WarnUsers  string    `json:"warn_users" grom:"comment:引用好友"`
	Content    string    `json:"content" grom:"comment:文本消息"`
	CreatedAt  time.Time `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt  time.Time `json:"updated_at" grom:"comment:更新时间"`
}
