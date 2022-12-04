package repo

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Base struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func NewBase(db *gorm.DB, rds *redis.Client) *Base {
	return &Base{Db: db, Redis: rds}
}
