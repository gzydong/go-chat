package model

import "time"

// 用户好友关系表
type UsersFriends struct {
	Id        int       `gorm:"column:id" json:"id" form:"id"`                         // 关系ID
	UserId    int       `gorm:"column:user_id" json:"user_id" form:"user_id"`          // 用户id
	FriendId  int       `gorm:"column:friend_id" json:"friend_id" form:"friend_id"`    // 好友id
	Remark    string    `gorm:"column:remark" json:"remark" form:"remark"`             // 好友的备注
	Status    int8      `gorm:"column:status" json:"status" form:"status"`             // 好友状态 [0:否;1:是]
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"` // 更新时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"` // 创建时间
}

type ContactListItem struct {
	Id       int    `gorm:"column:id" json:"id"`                // 用户ID
	Nickname string `gorm:"column:nickname" json:"nickname"`    // 用户昵称
	Gender   uint8  `gorm:"column:gender" json:"gender"`        // 用户性别[0:未知;1:男;2:女;]
	Motto    string `gorm:"column:motto" json:"motto"`          // 用户座右铭
	Avatar   string `grom:"column:avatar" json:"avatar" `       // 好友头像
	Remark   string `gorm:"column:remark" json:"friend_remark"` // 好友的备注
	IsOnline int    `json:"is_online"`                          // 是否在线
}
