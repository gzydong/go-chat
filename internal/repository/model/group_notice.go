package model

import (
	"database/sql"
	"time"
)

type GroupNotice struct {
	Id           int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 公告ID
	GroupId      int          `gorm:"column:group_id;default:0;NOT NULL" json:"group_id"`     // 群组ID
	CreatorId    int          `gorm:"column:creator_id;default:0;NOT NULL" json:"creator_id"` // 创建者用户ID
	Title        string       `gorm:"column:title;NOT NULL" json:"title"`                     // 公告标题
	Content      string       `gorm:"column:content;NOT NULL" json:"content"`                 // 公告内容
	ConfirmUsers string       `gorm:"column:confirm_users" json:"confirm_users"`              // 已确认成员
	IsDelete     int          `gorm:"column:is_delete;default:0;NOT NULL" json:"is_delete"`   // 是否删除[0:否;1:是;]
	IsTop        int          `gorm:"column:is_top;default:0;NOT NULL" json:"is_top"`         // 是否置顶[0:否;1:是;]
	IsConfirm    int          `gorm:"column:is_confirm;default:0;NOT NULL" json:"is_confirm"` // 是否需群成员确认公告[0:否;1:是;]
	CreatedAt    time.Time    `gorm:"column:created_at;NOT NULL" json:"created_at"`           // 创建时间
	UpdatedAt    time.Time    `gorm:"column:updated_at;NOT NULL" json:"updated_at"`           // 更新时间
	DeletedAt    sql.NullTime `gorm:"column:deleted_at" json:"deleted_at"`                    // 删除时间
}

func (GroupNotice) TableName() string {
	return "group_notice"
}

type SearchNoticeItem struct {
	Id           int       `json:"id" grom:"column:id"`
	CreatorId    int       `json:"creator_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	IsTop        int       `json:"is_top"`
	IsConfirm    int       `json:"is_confirm"`
	ConfirmUsers string    `json:"confirm_users"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Avatar       string    `json:"avatar"`
	Nickname     string    `json:"nickname"`
}
