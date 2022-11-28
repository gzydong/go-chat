package model

import "time"

type Emoticon struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 表情分组ID
	Name      string    `gorm:"column:name;NOT NULL" json:"name"`               // 分组名称
	Icon      string    `gorm:"column:icon;NOT NULL" json:"icon"`               // 分组图标
	Status    int       `gorm:"column:status;default:0;NOT NULL" json:"status"` // 分组状态[-1:已删除;0:正常;1:已禁用;]
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`   // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`   // 更新时间
}

func (Emoticon) TableName() string {
	return "emoticon"
}
