package service

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type BaseService struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewBaseService(db *gorm.DB, rds *redis.Client) *BaseService {
	return &BaseService{db, rds}
}

func (base *BaseService) Db() *gorm.DB {
	return base.db
}
