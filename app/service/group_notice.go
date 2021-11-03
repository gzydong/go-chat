package service

import (
	"context"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/timeutil"
	"time"
)

type GroupNoticeService struct {
	dao *dao.GroupNoticeDao
}

func NewGroupNoticeService(dao *dao.GroupNoticeDao) *GroupNoticeService {
	return &GroupNoticeService{
		dao: dao,
	}
}

func (s *GroupNoticeService) Dao() {

}

// Create 创建群公告
func (s *GroupNoticeService) Create(ctx context.Context, input *request.GroupNoticeEditRequest, uid int) error {
	notice := &model.GroupNotice{
		GroupId:   input.GroupId,
		CreatorId: uid,
		Title:     input.Title,
		Content:   input.Content,
		IsTop:     input.IsTop,
		IsConfirm: input.IsConfirm,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.dao.Db.Omit("deleted_at", "confirm_users").Create(notice).Error
}

// Update 更新群公告
func (s *GroupNoticeService) Update(ctx context.Context, input *request.GroupNoticeEditRequest) error {
	_, err := s.dao.Update(&model.GroupNotice{}, map[string]interface{}{
		"id":       input.NoticeId,
		"group_id": input.GroupId,
	}, map[string]interface{}{
		"title":      input.Title,
		"content":    input.Content,
		"is_top":     input.IsTop,
		"is_confirm": input.IsConfirm,
		"updated_at": time.Now(),
	})

	return err
}

func (s *GroupNoticeService) Delete(ctx context.Context, groupId, noticeId int) error {
	_, err := s.dao.Update(&model.GroupNotice{}, map[string]interface{}{
		"id":       noticeId,
		"group_id": groupId,
	}, map[string]interface{}{
		"is_delete":  1,
		"deleted_at": timeutil.DateTime(),
	})

	return err
}

func (s *GroupNoticeService) List(ctx context.Context, groupId int) []*model.SearchNoticeItem {

	items, _ := s.dao.GetListAll(ctx, groupId)

	return items
}
