package service

import (
	"context"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/timeutil"
	"gorm.io/gorm"
	"time"
)

type GroupNoticeService struct {
	db *gorm.DB
}

type NoticeItem struct {
	Id           int    `json:"id"`
	CreatorId    int    `json:"creator_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	IsTop        int    `json:"is_top"`
	IsConfirm    int    `json:"is_confirm"`
	ConfirmUsers string `json:"confirm_users"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
}

func NewGroupNoticeService(db *gorm.DB) *GroupNoticeService {
	return &GroupNoticeService{
		db: db,
	}
}

// Create 创建群公告
func (s *GroupNoticeService) Create(ctx context.Context, input *request.GroupNoticeEditRequest, userId int) error {
	notice := &model.GroupNotice{
		GroupId:   input.GroupId,
		CreatorId: userId,
		Title:     input.Title,
		Content:   input.Content,
		IsTop:     input.IsTop,
		IsConfirm: input.IsConfirm,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.db.Omit("deleted_at", "confirm_users").Create(notice).Error
}

func (s *GroupNoticeService) Update(ctx context.Context, input *request.GroupNoticeEditRequest, userId int) error {
	return s.db.Model(&model.GroupNotice{}).
		Select("title", "content", "updated_at", "is_top", "is_confirm").
		Where("id = ? and group_id = ?", input.NoticeId, input.GroupId).
		Updates(model.GroupNotice{
			Title:     input.Title,
			Content:   input.Content,
			IsTop:     input.IsTop,
			IsConfirm: input.IsConfirm,
			UpdatedAt: time.Now(),
		}).Error
}

func (s *GroupNoticeService) Delete(ctx context.Context, groupId, noticeId int) error {
	return s.db.Model(&model.GroupNotice{ID: noticeId, GroupId: groupId}).Updates(model.GroupNotice{
		IsDelete:  1,
		DeletedAt: timeutil.DateTime(),
	}).Error
}

func (s *GroupNoticeService) List(ctx context.Context, groupId int) []*NoticeItem {
	var items []*NoticeItem

	fields := []string{
		"lar_group_notice.id",
		"lar_group_notice.creator_id",
		"lar_group_notice.title",
		"lar_group_notice.content",
		"lar_group_notice.is_top",
		"lar_group_notice.is_confirm",
		"lar_group_notice.confirm_users",
		"lar_group_notice.created_at",
		"lar_group_notice.updated_at",
		"lar_users.avatar",
		"lar_users.nickname",
	}

	s.db.Table("lar_group_notice").
		Select(fields).
		Joins("left join lar_users on lar_users.id = lar_group_notice.creator_id").
		Where("lar_group_notice.group_id = ? and lar_group_notice.is_delete = ?", groupId, 0).
		Order("lar_group_notice.is_top desc").
		Order("lar_group_notice.updated_at desc").
		Scan(&items)

	return items
}
