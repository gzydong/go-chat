package organize

import (
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
)

type PositionService struct {
	*repo.Source
	dao *organize.Position
}

func NewPositionService(source *repo.Source, dao *organize.Position) *PositionService {
	return &PositionService{Source: source, dao: dao}
}
