package dao

import (
	"github.com/go-redis/redis/v8"
	"go-chat/internal/entity"
	"gorm.io/gorm"
)

type IBaseDao interface {
	BaseUpdate(model interface{}, where entity.MapStrAny, data entity.MapStrAny) (int, error)
}

type BaseDao struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewBaseDao(db *gorm.DB, rds *redis.Client) *BaseDao {
	return &BaseDao{db: db, rds: rds}
}

func (dao *BaseDao) Db() *gorm.DB {
	return dao.db
}

// BaseUpdate 批量更新
func (dao *BaseDao) BaseUpdate(model interface{}, where entity.MapStrAny, data entity.MapStrAny) (int, error) {
	fields := make([]string, 0, len(data))
	values := make(map[string]interface{})

	// 获取需要更新的字段
	for field, value := range data {
		fields = append(fields, field)
		values[field] = value
	}

	tx := dao.db.Model(model).Select(fields)

	for key, val := range where {
		tx.Where(key, val)
	}

	result := tx.Unscoped().Updates(values)

	return int(result.RowsAffected), result.Error
}
