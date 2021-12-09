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

func (s *TalkListService) GetTalkList(ctx context.Context, uid int) ([]*model.SearchTalkList, error) {
	var (
		err   error
		items = make([]*model.SearchTalkList, 0)
	)

	fields := []string{
		"list.id", "list.talk_type", "list.receiver_id", "list.updated_at",
		"list.is_disturb", "list.is_top", "list.is_robot",
		"`users`.avatar as user_avatar", "`users`.nickname",
		"`group`.group_name", "`group`.avatar as group_avatar",
	}

	tx := s.db.Table("talk_list list")
	tx.Select(fields)
	tx.Joins("left join `users` ON list.receiver_id = `users`.id AND list.talk_type = 1")
	tx.Joins("left join `group` ON list.receiver_id = `group`.id AND list.talk_type = 2")
	tx.Where("list.user_id = ? and list.is_delete = 0", uid)
	tx.Order("list.updated_at desc")

	if err = tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// Create 创建会话列表
func (s *TalkListService) Create(ctx context.Context, uid int, params *request.TalkListCreateRequest) (*model.TalkList, error) {
	var (
		err    error
		result *model.TalkList
	)

	err = s.db.Where(&model.TalkList{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}).First(&result).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		result = &model.TalkList{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		s.db.Create(result)
	} else {
		result.IsTop = 0
		result.IsDelete = 0
		result.IsDisturb = 0
		s.db.Save(result)
	}

	return result, nil
}

// Delete 删除会话
func (s *TalkListService) Delete(ctx context.Context, uid int, id int) error {
	return s.db.Model(&model.TalkList{}).Where("id = ? and user_id = ?", id, uid).Updates(map[string]interface{}{
		"is_delete":  1,
		"updated_at": time.Now(),
	}).Error
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
