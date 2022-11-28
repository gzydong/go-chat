package model

import "time"

// ContactApply 用户添加好友申请表
type ContactApply struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`       // 申请ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`     // 申请人ID
	FriendId  int       `gorm:"column:friend_id;default:0;NOT NULL" json:"friend_id"` // 被申请人
	Remark    string    `gorm:"column:remark;NOT NULL" json:"remark"`                 // 申请备注
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`         // 申请时间
}

func (ContactApply) TableName() string {
	return "contact_apply"
}

// ApplyItem 用户添加好友申请表
type ApplyItem struct {
	Id        int       `gorm:"column:id" json:"id"`                 // 申请ID
	UserId    int       `gorm:"column:user_id" json:"user_id"`       // 申请人ID
	FriendId  int       `gorm:"column:friend_id" json:"friend_id"`   // 被申请人
	Remark    string    `gorm:"column:remark" json:"remark"`         // 申请备注
	Nickname  string    `gorm:"column:nickname" json:"nickname"`     // 申请备注
	Avatar    string    `gorm:"column:avatar" json:"avatar"`         // 申请备注
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 申请时间
}
