package model

import "time"

// ContactGroup 联系人分组
type ContactGroup struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 主键ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户ID
	Name      string    `gorm:"column:name;NOT NULL" json:"remark"`               // 分组名称
	Num       int       `gorm:"column:num;default:0;NOT NULL" json:"num"`         // 成员总数
	Sort      int       `gorm:"column:sort;default:0;NOT NULL" json:"sort"`       // 分组名称
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`     // 更新时间
}

func (ContactGroup) TableName() string {
	return "contact_group"
}
