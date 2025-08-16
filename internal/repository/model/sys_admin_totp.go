package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type SysAdminTotp struct {
	Id          int         `gorm:"column:id" db:"id" json:"id"`
	AdminId     int         `gorm:"column:admin_id" db:"admin_id" json:"admin_id"`
	Secret      string      `gorm:"column:secret" db:"secret" json:"secret"`                                // 秘钥加密存储
	OneTimeCode StringSlice `gorm:"column:one_time_code;type:json" db:"one_time_code" json:"one_time_code"` // 一次性临时密码
	IsEnabled   string      `gorm:"column:is_enabled" db:"is_enabled" json:"is_enabled"`                    // 是否开启 Y:开启 N:关闭
	CreatedAt   time.Time   `gorm:"column:created_at;" db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;" db:"updated_at" json:"updated_at"`
}

func (SysAdminTotp) TableName() string {
	return "sys_admin_totp"
}

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("cannot scan %T into StringSlice", value)
	}
}
