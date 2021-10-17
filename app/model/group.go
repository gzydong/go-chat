package model

import (
	"database/sql"
)

type Group struct {
	ID          int          `json:"id" grom:"comment:群ID"`
	CreatorId   int          `json:"creator_id" grom:"comment:群主ID"`
	GroupName   string       `json:"group_name" grom:"comment:群名称"`
	Profile     string       `json:"profile" grom:"comment:群简介"`
	Avatar      string       `json:"avatar" grom:"comment:群头像"`
	MaxNum      int          `json:"max_num" grom:"comment:最大群成员数量"`
	IsOvert     int          `json:"is_overt" grom:"comment:是否公开可见"`
	IsMute      int          `json:"is_mute" grom:"comment:是否全员禁言"`
	IsDismiss   int          `json:"is_dismiss" grom:"comment:是否已解散"`
	CreatedAt   string       `json:"created_at" grom:"comment:创建时间"`
	DismissedAt sql.NullTime `json:"dismissed_at" grom:"comment:解散时间"`
}

func (m *Group) TableName() string {
	return "lar_group"
}
