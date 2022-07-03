package organize

import (
	"go-chat/internal/repository/dao/organize"
	"go-chat/internal/service"
)

type DeptService struct {
	*service.BaseService
	dao *organize.DepartmentDao
}

func NewOrganizeDeptService(baseService *service.BaseService, dao *organize.DepartmentDao) *DeptService {
	return &DeptService{BaseService: baseService, dao: dao}
}

func (s *DeptService) Dao() organize.IDeptDao {
	return s.dao
}
