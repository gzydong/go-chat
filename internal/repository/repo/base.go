package repo

import (
	"github.com/go-redis/redis/v8"
	"go-chat/internal/entity"
	"gorm.io/gorm"
)

type IBase interface {
	BaseUpdate(model interface{}, where entity.MapStrAny, data entity.MapStrAny) (int, error)
}

type Base struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func NewBase(db *gorm.DB, rds *redis.Client) *Base {
	return &Base{Db: db, Redis: rds}
}

// BaseUpdate 批量更新
func (b *Base) BaseUpdate(model interface{}, where entity.MapStrAny, data entity.MapStrAny) (int, error) {
	fields := make([]string, 0, len(data))
	values := make(map[string]interface{})

	// 获取需要更新的字段
	for field, value := range data {
		fields = append(fields, field)
		values[field] = value
	}

	tx := b.Db.Model(model).Select(fields)

	for key, val := range where {
		tx.Where(key, val)
	}

	result := tx.Unscoped().Updates(values)

	return int(result.RowsAffected), result.Error
}
