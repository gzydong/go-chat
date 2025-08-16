package model

import (
	"time"
)

// SysResource 接口资源表
type SysResource struct {
	Id        int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                                     // 资源ID
	Name      string    `gorm:"column:name;type:varchar(20);not null;unique" json:"name"`                                                         // 资源名称
	Uri       string    `gorm:"column:uri;type:varchar(255);not null;unique" json:"uri"`                                                          // 接口地址
	Type      int32     `gorm:"column:type;type:varchar(255);not null" json:"type"`                                                               // 类型 1:后台接口 2:开放接口
	Status    int32     `gorm:"column:status;type:tinyint(3);not null;default:1" json:"status"`                                                   // 状态[1:正常;2:停用;]
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`                             // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

// TableName 表名
func (SysResource) TableName() string {
	return "sys_resource"
}

// 资源类型常量
const (
	ResourceTypeAdmin = 1 // 后台接口
	ResourceTypeOpen  = 2 // 开放接口
)

// 资源状态常量
const (
	ResourceStatusNormal   = int32(1) // 正常
	ResourceStatusDisabled = int32(2) // 停用
)
