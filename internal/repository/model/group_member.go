package model

import (
	"time"
)

const (
	GroupMemberQuitStatusYes = 1
	GroupMemberQuitStatusNo  = 0

	GroupMemberMuteStatusYes = 1
	GroupMemberMuteStatusNo  = 0
)

type GroupMember struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`     // 自增ID
	GroupId   int       `gorm:"column:group_id;default:0;NOT NULL" json:"group_id"` // 群组ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`   // 用户ID
	Leader    int       `gorm:"column:leader;default:0;NOT NULL" json:"leader"`     // 成员属性[0:普通成员;1:管理员;2:群主;]
	UserCard  string    `gorm:"column:user_card;NOT NULL" json:"user_card"`         // 群名片
	IsQuit    int       `gorm:"column:is_quit;default:0;NOT NULL" json:"is_quit"`   // 是否退群[0:否;1:是;]
	IsMute    int       `gorm:"column:is_mute;default:0;NOT NULL" json:"is_mute"`   // 是否禁言[0:否;1:是;]
	JoinTime  time.Time `gorm:"column:join_time;" json:"join_time"`                 // 入群时间
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`       // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`       // 更新时间
}

func (GroupMember) TableName() string {
	return "group_member"
}

type MemberItem struct {
	Id       string `json:"id"`
	UserId   int    `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
	Motto    string `json:"motto"`
	Leader   int    `json:"leader"`
	IsMute   int    `json:"is_mute"`
	UserCard string `json:"user_card"`
}
