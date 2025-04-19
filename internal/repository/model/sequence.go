package model

import (
	"time"
)

type Sequence struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 主键ID
	SeqType   int32     `gorm:"column:seq_type;" json:"seq_type"`               // 业务类型 1用户 2群
	SourceId  int32     `gorm:"column:source_id;" json:"source_id"`             // 来源ID  type=1:用户ID，type=2:群ID
	CurSeq    int64     `gorm:"column:cur_seq;" json:"cur_seq"`                 // 当前分配ID
	MaxSeq    int64     `gorm:"column:max_seq;" json:"max_seq"`                 // 可发放最大ID
	Step      int64     `gorm:"column:step;" json:"step"`                       // 可发放最大ID
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (Sequence) TableName() string {
	return "sequence"
}
