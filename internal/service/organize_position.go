package service

import (
	"go-chat/internal/repository/repo"
)

type PositionService struct {
	*repo.Source
	Repo *repo.Position
}

func NewPositionService(source *repo.Source, dao *repo.Position) *PositionService {
	return &PositionService{Source: source, Repo: dao}
}
