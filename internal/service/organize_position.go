package service

import (
	"go-chat/internal/repository/repo"
)

type PositionService struct {
	*repo.Source
	Repo *repo.Position
}
