package service

import (
	"go-chat/internal/repository/repo"
)

type OrganizeService struct {
	*repo.Source
	Repo *repo.Organize
}
