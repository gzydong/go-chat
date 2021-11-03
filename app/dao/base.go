package dao

import (
	"errors"
	"gorm.io/gorm"
)

type Base struct {
	Db *gorm.DB
}

func NewBaseDao(db *gorm.DB) *Base {
	return &Base{db}
}

// Update 批量更新
func (b *Base) Update(model interface{}, where map[string]interface{}, data map[string]interface{}) (int, error) {

	fields := make([]string, len(data))

	// 获取需要更新的字段
	for field, _ := range data {
		fields = append(fields, field)
	}

	sql := b.Db.Model(model).Select(fields)

	for key, val := range where {
		sql.Where(key, val)
	}

	result := sql.Updates(data)

	return int(result.RowsAffected), result.Error
}

// FindByIds 根据主键查询一条或多条数据
func (b *Base) FindByIds(model interface{}, ids []int, fields interface{}) (bool, error) {
	var err error

	if len(ids) == 1 {
		err = b.Db.First(model, ids[0]).Error
	} else {
		err = b.Db.Select(fields).Find(model, ids).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}
