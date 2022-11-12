package organize

import (
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type Department struct {
	*repo.Base
}

func NewDepartment(base *repo.Base) *Department {
	return &Department{Base: base}
}

type IDept interface {
	FindAll() ([]*model.OrganizeDept, error)
}

func (d *Department) FindAll() ([]*model.OrganizeDept, error) {

	items := make([]*model.OrganizeDept, 0)

	err := d.Db.Model(model.OrganizeDept{}).Where("is_deleted = 1").Order("parent_id asc,order_num asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
