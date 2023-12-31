package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Sequence struct {
	db    *gorm.DB
	cache *cache.Sequence
}

func NewSequence(db *gorm.DB, cache *cache.Sequence) *Sequence {
	return &Sequence{db: db, cache: cache}
}

func (s *Sequence) try(ctx context.Context, userId int, receiverId int) error {
	result := s.cache.Redis().TTL(ctx, s.cache.Name(userId, receiverId)).Val()

	// 当数据不存在时需要从数据库中加载
	if result == time.Duration(-2) {
		lockName := fmt.Sprintf("%s_lock", s.cache.Name(userId, receiverId))

		isTrue := s.cache.Redis().SetNX(ctx, lockName, 1, 10*time.Second).Val()
		if !isTrue {
			return errors.New("请求频繁")
		}

		defer s.cache.Redis().Del(ctx, lockName)

		tx := s.db.WithContext(ctx).Model(&model.TalkRecords{})

		// 检测UserId 是否被设置，未设置则代表群聊
		if userId == 0 {
			tx = tx.Where("receiver_id = ? and talk_type = ?", receiverId, entity.ChatGroupMode)
		} else {
			tx = tx.Where("user_id = ? and receiver_id = ?", userId, receiverId).Or("user_id = ? and receiver_id = ?", receiverId, userId)
		}

		var seq int64
		err := tx.Select("ifnull(max(sequence),0)").Scan(&seq).Error
		if err != nil {
			logger.Errorf("[Sequence Total] 加载异常 err: %s", err.Error())
			return err
		}

		if err := s.cache.Set(ctx, userId, receiverId, seq); err != nil {
			logger.Errorf("[Sequence Set] 加载异常 err: %s", err.Error())
			return err
		}
	} else if result < time.Hour {
		s.cache.Redis().Expire(ctx, s.cache.Name(userId, receiverId), 12*time.Hour)
	}

	return nil
}

// Get 获取会话间的时序ID
func (s *Sequence) Get(ctx context.Context, userId int, receiverId int) int64 {

	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, userId, receiverId)
	}); err != nil {
		log.Println("Sequence GetObject Err :", err.Error())
	}

	return s.cache.Get(ctx, userId, receiverId)
}

// BatchGet 批量获取会话间的时序ID
func (s *Sequence) BatchGet(ctx context.Context, userId int, receiverId int, num int64) []int64 {

	if err := utils.Retry(5, 100*time.Millisecond, func() error {
		return s.try(ctx, userId, receiverId)
	}); err != nil {
		log.Println("Sequence BatchGet Err :", err.Error())
	}

	return s.cache.BatchGet(ctx, userId, receiverId, num)
}
