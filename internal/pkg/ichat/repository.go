package ichat

import (
	"context"

	"gorm.io/gorm"
)

type ITable interface {
	TableName() string
}

type Repo[T ITable] struct {
	model T        // 数据表结构体模型
	Db    *gorm.DB // 数据库
}

func NewRepo[T ITable](db *gorm.DB) Repo[T] {
	return Repo[T]{Db: db}
}

func (r *Repo[T]) Model(ctx context.Context) *gorm.DB {
	return r.Db.WithContext(ctx).Model(r.model)
}

// FindById 根据主键查询单条记录
func (r *Repo[T]) FindById(ctx context.Context, id int) (*T, error) {

	var item *T
	err := r.Db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

// FindByIds 根据主键查询单条记录
func (r *Repo[T]) FindByIds(ctx context.Context, ids []int) ([]*T, error) {

	var items []*T
	err := r.Db.WithContext(ctx).Find(&items, ids).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindAll 查询多条数据
// 注:仅支持单表
func (r *Repo[T]) FindAll(ctx context.Context, arg ...func(*gorm.DB)) ([]*T, error) {

	bd := r.Model(ctx)
	for _, fn := range arg {
		fn(bd)
	}

	var items []*T
	if err := bd.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByWhere 根据条件查询一条数据
func (r *Repo[T]) FindByWhere(ctx context.Context, where string, args ...interface{}) (*T, error) {

	var item *T
	err := r.Db.WithContext(ctx).Where(where, args...).First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

// QueryCount 根据条件统计数据总数
func (r *Repo[T]) QueryCount(ctx context.Context, where string, args ...interface{}) (int64, error) {

	var count int64
	err := r.Model(ctx).Where(where, args...).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// QueryExist 根据条件查询数据是否存在
func (r *Repo[T]) QueryExist(ctx context.Context, where string, args ...interface{}) (bool, error) {

	var count int64
	err := r.Model(ctx).Select("1").Where(where, args...).Limit(1).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// UpdateById 根据主键ID更新
func (r *Repo[T]) UpdateById(ctx context.Context, id interface{}, data map[string]interface{}) (int64, error) {
	res := r.Model(ctx).Where("id = ?", id).Updates(data)
	return res.RowsAffected, res.Error
}

// UpdateWhere 批量更新
func (r *Repo[T]) UpdateWhere(ctx context.Context, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	res := r.Model(ctx).Where(where, args...).Updates(data)
	return res.RowsAffected, res.Error
}
