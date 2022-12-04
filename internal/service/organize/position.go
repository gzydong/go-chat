package organize

import (
	"go-chat/internal/repository/repo/organize"
	"go-chat/internal/service"
)

type PositionService struct {
	*service.BaseService
	dao *organize.Position
}

func NewPositionService(baseService *service.BaseService, dao *organize.Position) *PositionService {
	return &PositionService{BaseService: baseService, dao: dao}
}

func (s *PositionService) Dao() *organize.Position {
	return s.dao
}
