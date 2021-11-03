package model

import (
	"time"
)

type GroupNotice struct {
	ID           int       `json:"id" grom:"comment:群公告ID"`
	GroupId      int       `json:"group_id" grom:"comment:群组ID"`
	CreatorId    int       `json:"creator_id" grom:"comment:创建者用户ID"`
	Title        string    `json:"title" grom:"comment:公告标题"`
	Content      string    `json:"content" grom:"comment:公告内容"`
	IsTop        int       `json:"is_top" grom:"comment:是否置顶"`
	IsDelete     int       `json:"is_delete" grom:"comment:是否删除"`
	IsConfirm    int       `json:"is_confirm" grom:"comment:是否需群成员确认公告"`
	ConfirmUsers string    `json:"confirm_users" grom:"comment:已确认成员"`
	CreatedAt    time.Time `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" grom:"comment:更新时间"`
	DeletedAt    string    `json:"-" grom:"comment:删除时间;default:'2021-10-23 12:27:50'"`
}

type SearchNoticeItem struct {
	Id           int    `json:"id"`
	CreatorId    int    `json:"creator_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	IsTop        int    `json:"is_top"`
	IsConfirm    int    `json:"is_confirm"`
	ConfirmUsers string `json:"confirm_users"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
}
