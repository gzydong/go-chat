package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseDao struct {
	db *gorm.DB
}

func NewBaseDao(db *gorm.DB) *BaseDao {
	return &BaseDao{db}
}

func (dao *BaseDao) Db() *gorm.DB {
	return dao.db
}

// BaseUpdate 批量更新
func (dao *BaseDao) BaseUpdate(model interface{}, where gin.H, data gin.H) (int, error) {
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

// FindByIds 根据主键查询一条或多条数据
func (dao *BaseDao) FindByIds(model interface{}, ids []int, fields interface{}) (bool, error) {
	var err error

	if len(ids) == 1 {
		err = dao.db.First(model, ids[0]).Error
	} else {
		err = dao.db.Select(fields).Find(model, ids).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}
