package organize

import (
	"go-chat/internal/dao/organize"
	"go-chat/internal/service"
)

type OrganizeService struct {
	*service.BaseService
	dao *organize.OrganizeDao
}

func NewOrganizeService(baseService *service.BaseService, dao *organize.OrganizeDao) *OrganizeService {
	return &OrganizeService{BaseService: baseService, dao: dao}
}

func (o *OrganizeService) Dao() organize.IOrganizeDao {
	return o.dao
}
