package organize

import (
	"go-chat/internal/repository/repo/organize"
	"go-chat/internal/service"
)

type DeptService struct {
	*service.BaseService
	dao *organize.Department
}

func NewOrganizeDeptService(baseService *service.BaseService, dao *organize.Department) *DeptService {
	return &DeptService{BaseService: baseService, dao: dao}
}

func (s *DeptService) Dao() *organize.Department {
	return s.dao
}
