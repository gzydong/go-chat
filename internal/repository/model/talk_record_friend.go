package model

import "time"

type TalkRecordFriend struct {
	Id         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`             // 聊天记录ID
	MsgId      string    `gorm:"column:msg_id;NOT NULL" json:"msg_id"`                       // 消息ID
	Sequence   int64     `gorm:"column:sequence;NOT NULL" json:"sequence"`                   // 消息时序ID（消息排序）
	MsgType    int       `gorm:"column:msg_type;default:1;NOT NULL" json:"msg_type"`         // 消息类型[1:文本消息;2:文件消息;3:会话消息;4:代码消息;5:投票消息;6:群公告;7:好友申请;8:登录通知;9:入群消息/退群消息;]
	UserId     int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`           // 用户ID
	FriendId   int       `gorm:"column:friend_id;default:0;NOT NULL" json:"friend_id"`       // 好友ID
	FromUserId int       `gorm:"column:from_user_id;default:0;NOT NULL" json:"from_user_id"` // 消息发送者ID
	QuoteId    string    `gorm:"column:quote_id;NOT NULL" json:"quote_id"`                   // 引用消息ID
	IsRevoke   int       `gorm:"column:is_revoke;default:0;NOT NULL" json:"is_revoke"`       // 是否撤回[0:否;1:是;]
	Extra      string    `gorm:"column:extra;NOT NULL" json:"extra"`                         // 消息扩展字段
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`               // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`               // 更新时间
}

func (TalkRecordFriend) TableName() string {
	return "talk_record_friend"
}
