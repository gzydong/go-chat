package model

import "time"

type OrganizePost struct {
	PositionId int       `gorm:"column:position_id;primary_key;AUTO_INCREMENT" json:"position_id"` // 岗位ID
	PostCode   string    `gorm:"column:post_code;NOT NULL" json:"post_code"`                       // 岗位编码
	PostName   string    `gorm:"column:post_name;NOT NULL" json:"post_name"`                       // 岗位名称
	Sort       int       `gorm:"column:sort;default:0;NOT NULL" json:"sort"`                       // 显示顺序
	Status     int       `gorm:"column:status;default:1;NOT NULL" json:"status"`                   // 状态[1:正常;2:停用;]
	Remark     string    `gorm:"column:remark;NOT NULL" json:"remark"`                             // 备注
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`                     // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`                     // 更新时间
}

func (OrganizePost) TableName() string {
	return "organize_position"
}
