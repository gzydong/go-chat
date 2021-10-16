package repository

import (
	"errors"
	"gorm.io/gorm"
)

type Base struct {
	db *gorm.DB
}

// Create 根据 model 创建一条数据
func (b *Base) Create(model interface{}) error {
	return b.db.Create(model).Error
}

func (b *Base) Update() {

}

// FindByIds 根据主键查询一条或多条数据
func (b *Base) FindByIds(model interface{}, ids []int, fields interface{}) (bool, error) {
	var err error

	if len(ids) == 1 {
		err = b.db.First(model, ids[0]).Error
	} else {
		err = b.db.Select(fields).Find(model, ids).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}
