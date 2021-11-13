package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupMember struct {
	ID        int            `json:"id" grom:"comment:群成员ID"`
	GroupId   int            `json:"group_id" grom:"comment:群组ID"`
	UserId    int            `json:"user_id" grom:"comment:用户ID"`
	Leader    int            `json:"leader" grom:"comment:成员属性"`
	IsMute    int            `json:"is_mute" grom:"comment:是否禁言"`
	IsQuit    int            `json:"is_quit" grom:"comment:是否退群"`
	UserCard  string         `json:"user_card" grom:"comment:群名片"`
	CreatedAt time.Time      `json:"created_at" grom:"comment:入群时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" grom:"comment:退群时间,"`
}

func (m *GroupMember) TableName() string {
	return "group_member"
}
