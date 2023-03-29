package organize

import (
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
)

type DeptService struct {
	*repo.Source
	dao *organize.Department
}

func NewOrganizeDeptService(source *repo.Source, dao *organize.Department) *DeptService {
	return &DeptService{Source: source, dao: dao}
}

func (s *DeptService) Dao() *organize.Department {
	return s.dao
}
