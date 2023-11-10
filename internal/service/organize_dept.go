package service

import (
	"go-chat/internal/repository/repo"
)

type DeptService struct {
	*repo.Source
	Repo *repo.Department
}

func NewOrganizeDeptService(source *repo.Source, dao *repo.Department) *DeptService {
	return &DeptService{Source: source, Repo: dao}
}
