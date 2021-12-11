package model

import (
	"database/sql"
	"time"
)

const (
	GroupMemberMaxNum = 200 // 最大成员数量
)

type Group struct {
	Id          int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 群ID
	CreatorId   int          `gorm:"column:creator_id;default:0" json:"creator_id"`  // 创建者ID(群主ID)
	GroupName   string       `gorm:"column:group_name" json:"group_name"`            // 群名称
	Profile     string       `gorm:"column:profile" json:"profile"`                  // 群介绍
	IsDismiss   int          `gorm:"column:is_dismiss;default:0" json:"is_dismiss"`  // 是否已解散[0:否;1:是;]
	Avatar      string       `gorm:"column:avatar" json:"avatar"`                    // 群头像
	MaxNum      int          `gorm:"column:max_num;default:200" json:"max_num"`      // 最大群成员数量
	IsOvert     int          `gorm:"column:is_overt;default:0" json:"is_overt"`      // 是否公开可见[0:否;1:是;]
	IsMute      int          `gorm:"column:is_mute;default:0" json:"is_mute"`        // 是否全员禁言 [0:否;1:是;]，提示:不包含群主或管理员
	CreatedAt   time.Time    `gorm:"column:created_at" json:"created_at"`            // 创建时间
	DismissedAt sql.NullTime `gorm:"column:dismissed_at" json:"dismissed_at"`        // 解散时间
}

func (m *Group) TableName() string {
	return "group"
}

type GroupItem struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
	Profile   string `json:"profile"`
	Leader    int    `json:"leader"`
	IsDisturb int    `json:"is_disturb"`
}
