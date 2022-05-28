package organize

import (
	"go-chat/internal/dao"
	"go-chat/internal/model"
)

type DepartmentDao struct {
	*dao.BaseDao
}

func NewDepartmentDao(baseDao *dao.BaseDao) *DepartmentDao {
	return &DepartmentDao{BaseDao: baseDao}
}

type IDeptDao interface {
	FindAll() ([]*model.OrganizeDept, error)
}

func (dao *DepartmentDao) FindAll() ([]*model.OrganizeDept, error) {

	items := make([]*model.OrganizeDept, 0)

	err := dao.Db().Model(model.OrganizeDept{}).Where("is_deleted = 1").Order("parent_id asc,order_num asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
