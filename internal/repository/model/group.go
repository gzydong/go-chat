package model

import (
	"time"
)

const (
	GroupMemberMaxNum = 200 // 最大成员数量
)

type Group struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 群ID
	Type      int       `gorm:"column:type;default:1;NOT NULL" json:"type"`             // 群类型[1:普通群;2:企业群;]
	CreatorId int       `gorm:"column:creator_id;default:0;NOT NULL" json:"creator_id"` // 创建者ID(群主ID)
	Name      string    `gorm:"column:name;NOT NULL" json:"name"`                       // 群名称
	Profile   string    `gorm:"column:profile;NOT NULL" json:"profile"`                 // 群介绍
	IsDismiss int       `gorm:"column:is_dismiss;default:0;NOT NULL" json:"is_dismiss"` // 是否已解散[0:否;1:是;]
	Avatar    string    `gorm:"column:avatar;NOT NULL" json:"avatar"`                   // 群头像
	MaxNum    int       `gorm:"column:max_num;default:200;NOT NULL" json:"max_num"`     // 最大群成员数量
	IsOvert   int       `gorm:"column:is_overt;default:0;NOT NULL" json:"is_overt"`     // 是否公开可见[0:否;1:是;]
	IsMute    int       `gorm:"column:is_mute;default:0;NOT NULL" json:"is_mute"`       // 是否全员禁言 [0:否;1:是;]，提示:不包含群主或管理员
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`           // 更新时间
}

func (Group) TableName() string {
	return "group"
}

type GroupItem struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
	Profile   string `json:"profile"`
	Leader    int    `json:"leader"`
	IsDisturb int    `json:"is_disturb"`
	CreatorId int    `json:"creator_id"`
}
