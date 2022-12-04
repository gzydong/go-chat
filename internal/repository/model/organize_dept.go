package model

import "time"

type OrganizeDept struct {
	DeptId    int       `gorm:"column:dept_id;primary_key;AUTO_INCREMENT" json:"dept_id"` // 部门id
	ParentId  int       `gorm:"column:parent_id;default:0;NOT NULL" json:"parent_id"`     // 父部门id
	Ancestors string    `gorm:"column:ancestors;NOT NULL" json:"ancestors"`               // 祖级列表
	DeptName  string    `gorm:"column:dept_name;NOT NULL" json:"dept_name"`               // 部门名称
	OrderNum  int       `gorm:"column:order_num;default:0;NOT NULL" json:"order_num"`     // 显示顺序
	Leader    string    `gorm:"column:leader;NOT NULL" json:"leader"`                     // 负责人
	Phone     string    `gorm:"column:phone;NOT NULL" json:"phone"`                       // 联系电话
	Email     string    `gorm:"column:email;NOT NULL" json:"email"`                       // 邮箱
	Status    int       `gorm:"column:status;default:1;NOT NULL" json:"status"`           // 部门状态[1:正常;2:停用]
	IsDeleted int       `gorm:"column:is_deleted;default:1;NOT NULL" json:"is_deleted"`   // 是否删除[1:否;2:是]
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
}

func (OrganizeDept) TableName() string {
	return "organize_dept"
}
