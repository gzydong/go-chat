package repo

import (
	"context"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
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

func (s *Sequence) try(ctx context.Context, userId int, receiverId int) {
	result := s.cache.Redis().TTL(ctx, s.cache.Name(userId, receiverId)).Val()

	// 当数据不存在时需要从数据库中加载
	// 这里可能存在并发问题，但会话间 Sequence ID 并发情况下从复也几乎是能忍受的
	if result == time.Duration(-2) {
		tx := s.db.WithContext(ctx).Model(&model.TalkRecords{})

		// 检测UserId 是否被设置，未设置则代表群聊
		if userId == 0 {
			tx = tx.Where("receiver_id = ? and type = ?", receiverId, entity.ChatGroupMode)
		} else {
			tx = tx.Where("user_id = ? and receiver_id = ?", userId, receiverId).Or("user_id = ? and receiver_id = ?", receiverId, userId)
		}

		var seq int64
		err := tx.Select("max(sequence)").Scan(&seq).Error
		if err != nil {
			logger.Error("[Sequence Total] 加载异常 err: ", err.Error())
			return
		}

		if err := s.cache.Init(ctx, userId, receiverId, seq); err != nil {
			logger.Error("[Sequence Init] 加载异常 err: ", err.Error())
		}
	} else if result == time.Duration(-1) {
		s.cache.Redis().Expire(ctx, s.cache.Name(userId, receiverId), 12*time.Hour)
	}
}

// Get 获取会话间的时序ID
func (s *Sequence) Get(ctx context.Context, userId int, receiverId int) int64 {

	s.try(ctx, userId, receiverId)

	return s.cache.Get(ctx, userId, receiverId)
}

// BatchGet 批量获取会话间的时序ID
func (s *Sequence) BatchGet(ctx context.Context, userId int, receiverId int, num int64) []int64 {

	s.try(ctx, userId, receiverId)

	return s.cache.BatchGet(ctx, userId, receiverId, num)
}
