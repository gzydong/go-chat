package repo

import (
	"context"
	"log"
	"time"

	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SequenceType int32

const (
	SequenceTypeUser  = 1
	SequenceTypeGroup = 2
)

type Sequence struct {
	db    *gorm.DB
	cache *cache.Sequence
	core.Repo[model.Sequence]
}

func NewSequence(db *gorm.DB, cache *cache.Sequence) *Sequence {
	return &Sequence{db: db, cache: cache, Repo: core.NewRepo[model.Sequence](db)}
}

func (s *Sequence) try(ctx context.Context, seqType SequenceType, sourceId int32) error {
	key := s.cache.Name(int32(seqType), sourceId)

	// 当数据不存在时需要从数据库中加载
	if s.cache.Redis().TTL(ctx, key).Val() == time.Duration(-2) {
		// TODO 没有缓存时，改怎么处理
		_ = s.cache.Set(ctx, int32(seqType), sourceId, 100000)
	}

	return nil
}

// Get 获取会话间的时序ID
func (s *Sequence) Get(ctx context.Context, seqType SequenceType, sourceId int32) int64 {
	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, seqType, sourceId)
	}); err != nil {
		log.Println("Sequence Get Err :", err.Error())
	}

	return s.cache.Get(ctx, int32(seqType), sourceId)
}

// BatchGet 批量获取会话间的时序ID
func (s *Sequence) BatchGet(ctx context.Context, seqType SequenceType, sourceId int32, num int) []int64 {
	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, seqType, sourceId)
	}); err != nil {
		log.Println("Sequence BatchGet Err :", err.Error())
	}

	return s.cache.BatchGet(ctx, int32(seqType), sourceId, int64(num))
}
