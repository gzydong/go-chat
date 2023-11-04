package service

import (
	"go-chat/internal/repository/repo"
)

type OrganizeService struct {
	*repo.Source
	Repo *repo.Organize
}

func NewOrganizeService(source *repo.Source, dao *repo.Organize) *OrganizeService {
	return &OrganizeService{Source: source, Repo: dao}
}
