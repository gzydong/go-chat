package model

import "time"

// ContactGroup 联系人分组
type ContactGroup struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 主键ID
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	Name      string    `gorm:"column:name;" json:"remark"`                     // 分组名称
	Num       int       `gorm:"column:num;" json:"num"`                         // 成员总数
	Sort      int       `gorm:"column:sort;" json:"sort"`                       // 分组名称
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (ContactGroup) TableName() string {
	return "contact_group"
}
