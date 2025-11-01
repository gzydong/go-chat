package service

import (
	"github.com/gzydong/go-chat/internal/repository/repo"
)

type OrganizeService struct {
	*repo.Source
	Repo *repo.Organize
}
