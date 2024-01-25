package model

import (
	"time"
)

const (
	GroupMemberLeaderOwner    = 1 // 群主
	GroupMemberLeaderAdmin    = 2 // 管理员
	GroupMemberLeaderOrdinary = 3 // 普通成员
)

type GroupMember struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 自增ID
	GroupId   int       `gorm:"column:group_id;" json:"group_id"`               // 群组ID
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	Leader    int       `gorm:"column:leader;" json:"leader"`                   // 成员属性[1:群主;2:管理员;3:普通成员;]
	UserCard  string    `gorm:"column:user_card;" json:"user_card"`             // 群名片
	IsQuit    int       `gorm:"column:is_quit;" json:"is_quit"`                 // 是否退群[1:否;2:是;]
	IsMute    int       `gorm:"column:is_mute;" json:"is_mute"`                 // 是否禁言[1:否;2:是;]
	JoinTime  time.Time `gorm:"column:join_time;" json:"join_time"`             // 入群时间
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
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
