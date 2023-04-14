package model

import "time"

const (
	ContactStatusNormal = 1
	ContactStatusDelete = 0
)

// Contact 用户好友关系表
type Contact struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                         // 关系ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`                       // 用户id
	FriendId  int       `gorm:"column:friend_id;default:0;NOT NULL" json:"friend_id"`                   // 好友id
	Remark    string    `gorm:"column:remark;NOT NULL" json:"remark"`                                   // 好友的备注
	Status    int       `gorm:"column:status;default:0;NOT NULL" json:"status"`                         // 好友状态 [0:否;1:是]
	GroupId   int       `gorm:"column:group_id;default:0;NOT NULL" json:"group_id"`                     // 分组id
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 更新时间
}

func (Contact) TableName() string {
	return "contact"
}

type ContactListItem struct {
	Id       int    `gorm:"column:id" json:"id"`                // 用户ID
	Nickname string `gorm:"column:nickname" json:"nickname"`    // 用户昵称
	Gender   uint8  `gorm:"column:gender" json:"gender"`        // 用户性别[0:未知;1:男;2:女;]
	Motto    string `gorm:"column:motto" json:"motto"`          // 用户座右铭
	Avatar   string `grom:"column:avatar" json:"avatar" `       // 好友头像
	Remark   string `gorm:"column:remark" json:"friend_remark"` // 好友的备注
	IsOnline int    `json:"is_online"`                          // 是否在线
	GroupId  int    `gorm:"column:group_id" json:"group_id"`    // 联系人分组
}
