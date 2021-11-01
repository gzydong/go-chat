package service

import (
	"context"
	"errors"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"gorm.io/gorm"
	"time"
)

type TalkListService struct {
	*BaseService
	dao *dao.TalkListDao
}

func NewTalkListService(base *BaseService, dao *dao.TalkListDao) *TalkListService {
	return &TalkListService{base, dao}
}

func (s *TalkListService) Dao() *dao.TalkListDao {
	return s.dao
}

// Create 创建会话列表
func (s *TalkListService) Create(ctx context.Context, uid int, params *request.TalkListCreateRequest) (*model.TalkList, error) {
	var (
		err    error
		result model.TalkList
	)

	err = s.db.Debug().Where(&model.TalkList{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		result = model.TalkList{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		s.db.Create(&result)
	} else {
		result.IsTop = 0
		result.IsDelete = 0
		result.IsDisturb = 0
		s.db.Save(result)
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
