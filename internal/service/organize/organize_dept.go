package organize

import (
	"go-chat/internal/dao/organize"
	"go-chat/internal/service"
)

type OrganizeDeptService struct {
	*service.BaseService
	dao *organize.DepartmentDao
}

func NewOrganizeDeptService(baseService *service.BaseService, dao *organize.DepartmentDao) *OrganizeDeptService {
	return &OrganizeDeptService{BaseService: baseService, dao: dao}
}

func (s *OrganizeDeptService) Dao() organize.IDeptDao {
	return s.dao
}
