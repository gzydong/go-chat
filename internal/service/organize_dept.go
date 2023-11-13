package service

import (
	"go-chat/internal/repository/repo"
)

type DeptService struct {
	*repo.Source
	Repo *repo.Department
}
