package organize

import (
	"go-chat/internal/repository/repo/organize"
	"go-chat/internal/service"
)

type OrganizeService struct {
	*service.BaseService
	dao *organize.Organize
}

func NewOrganizeService(baseService *service.BaseService, dao *organize.Organize) *OrganizeService {
	return &OrganizeService{BaseService: baseService, dao: dao}
}

func (o *OrganizeService) Dao() *organize.Organize {
	return o.dao
}
