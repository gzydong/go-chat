package model

type Emoticon struct {
	Id        int      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 表情分组ID
	Name      string   `gorm:"column:name;NOT NULL" json:"name"`               // 分组名称
	Icon      string   `gorm:"column:icon" json:"icon"`                        // 分组图标
	CreatedAt int64    `gorm:"column:created_at;default:0" json:"created_at"`  // 创建时间
	Status    int      `gorm:"column:status;default:0" json:"status"`          // 分组状态[-1:已删除;0:正常;1:已禁用;]
	UpdatedAt DateTime `gorm:"column:updated_at" json:"updated_at"`            // 更新时间
}
