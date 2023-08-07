package organize

import (
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
)

type OrganizeService struct {
	*repo.Source
	dao *organize.Organize
}

func NewOrganizeService(source *repo.Source, dao *organize.Organize) *OrganizeService {
	return &OrganizeService{Source: source, dao: dao}
}
