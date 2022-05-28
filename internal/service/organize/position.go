package organize

import (
	"go-chat/internal/dao/organize"
	"go-chat/internal/service"
)

type PositionService struct {
	*service.BaseService
	dao *organize.PositionDao
}

func NewPositionService(baseService *service.BaseService, dao *organize.PositionDao) *PositionService {
	return &PositionService{BaseService: baseService, dao: dao}
}

func (s *PositionService) Dao() organize.IPositionDao {
	return s.dao
}
