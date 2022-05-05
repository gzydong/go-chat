package model

import "time"

type TalkRecords struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`           // 聊天记录ID
	TalkType   int       `gorm:"column:talk_type;default:1;NOT NULL" json:"talk_type"`     // 对话类型[1:私信;2:群聊;]
	MsgType    int       `gorm:"column:msg_type;default:0;NOT NULL" json:"msg_type"`       // 消息类型[0:系统消息;1:文本消息;2:文件消息;3:会话消息;4:代码消息;5:投票消息;6:群公告;7:好友申请;8:登录通知;9:入群消息/退群消息;]
	UserId     int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`         // 发送者ID（用户ID）
	ReceiverId int       `gorm:"column:receiver_id;default:0;NOT NULL" json:"receiver_id"` // 接收者ID（用户ID 或 群ID）
	IsRevoke   int       `gorm:"column:is_revoke;default:0;NOT NULL" json:"is_revoke"`     // 是否撤回消息[0:否;1:是;]
	IsMark     int       `gorm:"column:is_mark;default:0;NOT NULL" json:"is_mark"`         // 是否重要消息[0:否;1:是;]
	IsRead     int       `gorm:"column:is_read;default:0;NOT NULL" json:"is_read"`         // 是否已读[0:否;1:是;]
	QuoteId    int       `gorm:"column:quote_id;default:0;NOT NULL" json:"quote_id"`       // 引用消息ID
	WarnUsers  string    `gorm:"column:warn_users;NOT NULL" json:"warn_users"`             // @好友 、 多个用英文逗号 “,” 拼接 (0:代表所有人)
	Content    string    `gorm:"column:content" json:"content"`                            // 文本消息 {@nickname@}
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
}

type QueryTalkRecordsItem struct {
	Id         int       `json:"id"`
	TalkType   int       `json:"talk_type"`
	MsgType    int       `json:"msg_type"`
	UserId     int       `json:"user_id"`
	ReceiverId int       `json:"receiver_id"`
	IsRevoke   int       `json:"is_revoke"`
	IsMark     int       `json:"is_mark"`
	IsRead     int       `json:"is_read"`
	QuoteId    int       `json:"quote_id"`
	WarnUsers  string    `json:"warn_users"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
}
