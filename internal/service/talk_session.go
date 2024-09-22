package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"go-chat/internal/entity"
	"strings"
	"time"

	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ ITalkSessionService = (*TalkSessionService)(nil)

type ITalkSessionService interface {
	List(ctx context.Context, uid int) ([]*model.SearchTalkSession, error)
	Create(ctx context.Context, opt *TalkSessionCreateOpt) (*model.TalkSession, error)
	Delete(ctx context.Context, uid int, talkMode int, toFromId int) error
	Top(ctx context.Context, opt *TalkSessionTopOpt) error
	Disturb(ctx context.Context, opt *TalkSessionDisturbOpt) error
	BatchAddList(ctx context.Context, uid int, values map[string]int)
}

type TalkSessionService struct {
	*repo.Source
	TalkSessionRepo *repo.TalkSession
}

func (s *TalkSessionService) List(ctx context.Context, uid int) ([]*model.SearchTalkSession, error) {

	fields := []string{
		"list.id", "list.talk_mode", "list.to_from_id", "list.updated_at",
		"list.is_disturb", "list.is_top", "list.is_robot",
		"`users`.avatar", "`users`.nickname",
		"`group`.name as group_name", "`group`.avatar as group_avatar",
	}

	query := s.Source.Db().WithContext(ctx).Table("talk_session list")
	query.Joins("left join `users` ON list.to_from_id = `users`.id AND list.talk_mode = ?", entity.ChatPrivateMode)
	query.Joins("left join `group` ON list.to_from_id = `group`.id AND list.talk_mode = ?", entity.ChatGroupMode)
	query.Where("list.user_id = ? and list.is_delete = ?", uid, model.No)
	query.Order("list.updated_at desc")

	var items []*model.SearchTalkSession
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

type TalkSessionCreateOpt struct {
	UserId     int
	TalkType   int
	ReceiverId int
	IsBoot     bool
}

// Create 创建会话列表
func (s *TalkSessionService) Create(ctx context.Context, opt *TalkSessionCreateOpt) (*model.TalkSession, error) {

	result, err := s.TalkSessionRepo.FindByWhere(ctx, "talk_mode = ? and user_id = ? and to_from_id = ?", opt.TalkType, opt.UserId, opt.ReceiverId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		result = &model.TalkSession{
			TalkMode:  opt.TalkType,
			UserId:    opt.UserId,
			ToFromId:  opt.ReceiverId,
			IsTop:     model.No,
			IsDelete:  model.No,
			IsDisturb: model.No,
			IsRobot:   model.No,
		}

		if opt.IsBoot {
			result.IsRobot = model.Yes
		}

		s.Source.Db().WithContext(ctx).Create(result)
	} else {
		result.IsTop = model.No
		result.IsDelete = model.No
		result.IsDisturb = model.No

		if opt.IsBoot {
			result.IsRobot = model.Yes
		}

		s.Source.Db().WithContext(ctx).Save(result)
	}

	return result, nil
}

// Delete 删除会话
func (s *TalkSessionService) Delete(ctx context.Context, uid int, talkMode int, toFromId int) error {
	_, err := s.TalkSessionRepo.UpdateByWhere(ctx, map[string]any{
		"is_delete":  model.Yes,
		"updated_at": time.Now(),
	}, "user_id = ? and to_from_id = ? and talk_mode = ?", uid, toFromId, talkMode)
	return err
}

type TalkSessionTopOpt struct {
	UserId   int // 用户id
	TalkMode int // 1:私聊 2:群聊
	ToFromId int // 对方id
	Action   int // 1:置顶 2:取消置顶
}

// Top 会话置顶
func (s *TalkSessionService) Top(ctx context.Context, opt *TalkSessionTopOpt) error {
	_, err := s.TalkSessionRepo.UpdateByWhere(ctx, map[string]any{
		"is_top":     lo.Ternary(opt.Action == 1, model.Yes, model.No),
		"updated_at": time.Now(),
	}, "user_id = ? and talk_mode = ? and to_from_id = ?", opt.UserId, opt.TalkMode, opt.ToFromId)
	return err
}

type TalkSessionDisturbOpt struct {
	UserId   int
	TalkMode int
	ToFromId int
	Action   int
}

// Disturb 会话免打扰
func (s *TalkSessionService) Disturb(ctx context.Context, opt *TalkSessionDisturbOpt) error {
	_, err := s.TalkSessionRepo.UpdateByWhere(ctx, map[string]any{
		"is_disturb": lo.Ternary(opt.Action == 1, model.Yes, model.No),
		"updated_at": time.Now(),
	}, "user_id = ? and talk_mode = ? and to_from_id = ?", opt.UserId, opt.TalkMode, opt.ToFromId)
	return err
}

// BatchAddList 批量添加会话列表
func (s *TalkSessionService) BatchAddList(ctx context.Context, uid int, values map[string]int) {

	ctime := timeutil.DateTime()

	data := make([]string, 0)
	for k, v := range values {
		if v == 0 {
			continue
		}

		value := strings.Split(k, "_")
		if len(value) != 2 {
			continue
		}

		data = append(data, fmt.Sprintf("(%s, %d, %s, '%s', '%s')", value[0], uid, value[1], ctime, ctime))
	}

	if len(data) == 0 {
		return
	}

	s.Source.Db().WithContext(ctx).Exec(fmt.Sprintf("INSERT INTO talk_session ( `talk_mode`, `user_id`, `to_from_id`, created_at, updated_at ) VALUES %s ON DUPLICATE KEY UPDATE is_delete = %d, updated_at = '%s'", strings.Join(data, ","), model.No, ctime))
}
