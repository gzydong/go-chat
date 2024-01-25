package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"gorm.io/gorm"
)

type Sequence struct {
	db    *gorm.DB
	cache *cache.Sequence
}

func NewSequence(db *gorm.DB, cache *cache.Sequence) *Sequence {
	return &Sequence{db: db, cache: cache}
}

func (s *Sequence) try(ctx context.Context, id int, isUserId bool) error {
	result := s.cache.Redis().TTL(ctx, s.cache.Name(id, isUserId)).Val()

	// 当数据不存在时需要从数据库中加载
	if result == time.Duration(-2) {
		lockName := fmt.Sprintf("%s_lock", s.cache.Name(id, isUserId))

		isTrue := s.cache.Redis().SetNX(ctx, lockName, 1, 10*time.Second).Val()
		if !isTrue {
			return errors.New("请求频繁")
		}

		defer s.cache.Redis().Del(ctx, lockName)

		tx := s.db.WithContext(ctx).Select("ifnull(max(sequence),0)")
		if isUserId {
			tx.Table("talk_user_message").Where("user_id = ?", id)
		} else {
			tx.Table("talk_group_message").Where("group_id = ?", id)
		}

		var seq int64
		if err := tx.Scan(&seq).Error; err != nil {
			logger.Errorf("[Sequence Total] 加载异常 err: %s", err.Error())
			return err
		}

		if err := s.cache.Set(ctx, id, isUserId, seq+100); err != nil {
			logger.Errorf("[Sequence set] 加载异常 err: %s", err.Error())
			return err
		}
	} else if result < time.Hour {
		s.cache.Redis().Expire(ctx, s.cache.Name(id, isUserId), 12*time.Hour)
	}

	return nil
}

// Get 获取会话间的时序ID
func (s *Sequence) Get(ctx context.Context, id int, isUserId bool) int64 {

	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, id, isUserId)
	}); err != nil {
		log.Println("Sequence GetObject Err :", err.Error())
	}

	return s.cache.Get(ctx, id, isUserId)
}

// BatchGet 批量获取会话间的时序ID
func (s *Sequence) BatchGet(ctx context.Context, id int, isUserId bool, num int64) []int64 {

	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, id, isUserId)
	}); err != nil {
		log.Println("Sequence BatchGet Err :", err.Error())
	}

	return s.cache.BatchGet(ctx, id, isUserId, num)
}
