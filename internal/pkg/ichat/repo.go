package ichat

import (
	"context"

	"gorm.io/gorm"
)

type ITable interface {
	TableName() string
}

type Base[T ITable] struct {
	model T        // 数据表结构体模型
	Db    *gorm.DB // 数据库
}

// FindById 根据主键查询单条记录
func (b *Base[T]) FindById(ctx context.Context, id int) (*T, error) {
	var data *T
	err := b.Db.WithContext(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindAll 查询多条数据
// 注:仅支持单表
func (b *Base[T]) FindAll(ctx context.Context, arg ...func(*gorm.DB)) ([]*T, error) {

	var items []*T
	bd := b.Db.WithContext(ctx).Model(b.model)

	for _, fn := range arg {
		fn(bd)
	}

	if err := bd.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByWhere 根据条件查询一条数据
func (b *Base[T]) FindByWhere(ctx context.Context, where string, args ...interface{}) (*T, error) {

	var data *T
	err := b.Db.WithContext(ctx).Where(where, args...).First(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindByCount 根据条件统计数据总数
func (b *Base[T]) FindByCount(ctx context.Context, where string, args ...interface{}) (int64, error) {
	var count int64

	err := b.Db.WithContext(ctx).Model(b.model).Where(where, args...).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateById 根据主键ID更新
func (b *Base[T]) UpdateById(ctx context.Context, id interface{}, data map[string]interface{}) (int64, error) {
	res := b.Db.Debug().WithContext(ctx).Model(b.model).Where("id = ?", id).Updates(data)

	return res.RowsAffected, res.Error
}

// Updates 批量更新
func (b *Base[T]) Updates(ctx context.Context, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	res := b.Db.Debug().WithContext(ctx).Model(b.model).Where(where, args...).Updates(data)

	return res.RowsAffected, res.Error
}
