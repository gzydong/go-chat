package service

import (
	"github.com/gzydong/go-chat/internal/repository/repo"
)

type DeptService struct {
	*repo.Source
	Repo *repo.Department
}
