package service

import (
	"context"

	"go-chat/internal/repository/cache"
)

var _ IRoomService = (*RoomService)(nil)

type RoomMemberOpt struct {
	Channel  string // 渠道分类
	RoomType int    // 房间类型
	Number   string // 房间号
	Sid      string // 服务ID
	Cid      int64  // 客户端ID
}

type RoomOpt struct {
	Channel  string // 渠道分类
	RoomType int    // 房间类型
	Number   string // 房间号
	Sid      string // 服务ID
}

type IRoomService interface {
	// AddMember 添加房间成员
	AddMember(ctx context.Context, opt RoomMemberOpt) error
	// MultiAddMember 批量添加房间成员
	MultiAddMember(ctx context.Context, opts []RoomMemberOpt) error
	// DelMember 删除房间成员
	DelMember(ctx context.Context, opt RoomMemberOpt) error
	// MultiDelMember 批量删除房间成员
	MultiDelMember(ctx context.Context, opts []RoomMemberOpt) error
	// FindAllClientIds 获取房间所有成员
	FindAllClientIds(ctx context.Context, opt RoomOpt) ([]int64, error)
}

type RoomService struct {
	RoomStorage *cache.RoomStorage
}

func (r *RoomService) AddMember(ctx context.Context, opt RoomMemberOpt) error {
	return r.RoomStorage.Add(ctx, &cache.RoomOption{
		Channel:  opt.Channel,
		RoomType: opt.RoomType,
		Number:   opt.Number,
		Sid:      opt.Sid,
		Cid:      opt.Cid,
	})
}

func (r *RoomService) MultiAddMember(ctx context.Context, opts []RoomMemberOpt) error {
	items := make([]*cache.RoomOption, 0, len(opts))
	for _, opt := range opts {
		items = append(items, &cache.RoomOption{
			Channel:  opt.Channel,
			RoomType: opt.RoomType,
			Number:   opt.Number,
			Sid:      opt.Sid,
			Cid:      opt.Cid,
		})
	}

	return r.RoomStorage.BatchAdd(ctx, items)
}

func (r *RoomService) DelMember(ctx context.Context, opt RoomMemberOpt) error {
	return r.RoomStorage.Del(ctx, &cache.RoomOption{
		Channel:  opt.Channel,
		RoomType: opt.RoomType,
		Number:   opt.Number,
		Sid:      opt.Sid,
		Cid:      opt.Cid,
	})
}

func (r *RoomService) MultiDelMember(ctx context.Context, opts []RoomMemberOpt) error {
	items := make([]*cache.RoomOption, 0, len(opts))
	for _, opt := range opts {
		items = append(items, &cache.RoomOption{
			Channel:  opt.Channel,
			RoomType: opt.RoomType,
			Number:   opt.Number,
			Sid:      opt.Sid,
			Cid:      opt.Cid,
		})
	}

	return r.RoomStorage.BatchDel(ctx, items)
}

func (r *RoomService) FindAllClientIds(ctx context.Context, opt RoomOpt) ([]int64, error) {
	return r.RoomStorage.All(ctx, &cache.RoomOption{
		Channel:  opt.Channel,
		RoomType: opt.RoomType,
		Number:   opt.Number,
		Sid:      opt.Sid,
	}), nil
}
