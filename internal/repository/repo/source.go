package repo

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Source struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewSource(db *gorm.DB, redis *redis.Client) *Source {
	return &Source{db: db, redis: redis}
}

func (s *Source) Db() *gorm.DB {
	return s.db
}

func (s *Source) Redis() *redis.Client {
	return s.redis
}
