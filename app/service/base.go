package service

import "gorm.io/gorm"

type BaseService struct {
	db *gorm.DB
}

func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{db}
}
