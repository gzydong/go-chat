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

// Delete 删除会话
func (s *TalkListService) Delete(ctx context.Context, uid int, id int) error {
	err := s.db.Model(&model.TalkList{}).Where("id = ? and user_id = ?", id, uid).Updates(map[string]interface{}{
		"is_delete":  1,
		"updated_at": time.Now(),
	}).Error

	return err
}

// Top 会话置顶
func (s *TalkListService) Top(ctx context.Context, uid int, params *request.TalkListTopRequest) error {

	isTop := 0

	if params.Type == 1 {
		isTop = 1
	}

	err := s.db.Model(&model.TalkList{}).Where("id = ? and user_id = ?", params.Id, uid).
		Updates(map[string]interface{}{
			"is_top":     isTop,
			"updated_at": time.Now(),
		}).Error

	return err
}

// Top 会话置顶
func (s *TalkListService) Disturb(ctx context.Context, uid int, params *request.TalkListDisturbRequest) error {
	err := s.db.Model(&model.TalkList{}).
		Where("user_id = ? and receiver_id = ? and talk_type = ?", uid, params.ReceiverId, params.TalkType).
		Updates(map[string]interface{}{
			"is_disturb": params.IsDisturb,
			"updated_at": time.Now(),
		}).Error

	return err
}
