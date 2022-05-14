package service

import (
	"context"
	"time"

	"go-chat/internal/dao"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/timeutil"
)

type GroupNoticeEditOpts struct {
	UserId    int
	GroupId   int
	NoticeId  int
	Title     string
	Content   string
	IsTop     int
	IsConfirm int
}

type GroupNoticeService struct {
	dao *dao.GroupNoticeDao
}

func NewGroupNoticeService(dao *dao.GroupNoticeDao) *GroupNoticeService {
	return &GroupNoticeService{
		dao: dao,
	}
}

func (s *GroupNoticeService) Dao() *dao.GroupNoticeDao {
	return s.dao
}

// Create 创建群公告
func (s *GroupNoticeService) Create(ctx context.Context, opts *GroupNoticeEditOpts) error {
	return s.dao.Db().Create(&model.GroupNotice{
		GroupId:      opts.GroupId,
		CreatorId:    opts.UserId,
		Title:        opts.Title,
		Content:      opts.Content,
		IsTop:        opts.IsTop,
		IsConfirm:    opts.IsConfirm,
		ConfirmUsers: "{}",
	}).Error
}

// Update 更新群公告
func (s *GroupNoticeService) Update(ctx context.Context, opts *GroupNoticeEditOpts) error {
	_, err := s.dao.BaseUpdate(&model.GroupNotice{}, entity.MapStrAny{
		"id":       opts.NoticeId,
		"group_id": opts.GroupId,
	}, entity.MapStrAny{
		"title":      opts.Title,
		"content":    opts.Content,
		"is_top":     opts.IsTop,
		"is_confirm": opts.IsConfirm,
		"updated_at": time.Now(),
	})

	return err
}

func (s *GroupNoticeService) Delete(ctx context.Context, groupId, noticeId int) error {
	_, err := s.dao.BaseUpdate(&model.GroupNotice{}, entity.MapStrAny{
		"id":       noticeId,
		"group_id": groupId,
	}, entity.MapStrAny{
		"is_delete":  1,
		"deleted_at": timeutil.DateTime(),
	})

	return err
}
