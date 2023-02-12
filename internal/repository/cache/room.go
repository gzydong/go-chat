package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/entity"
)

type RoomStorage struct {
	rds *redis.Client
}

type RoomOption struct {
	Channel  string          // 渠道分类
	RoomType entity.RoomType // 房间类型
	Number   string          // 房间号
	Sid      string          // 网关ID
	Cid      int64           // 客户端ID
}

func NewRoomStorage(rds *redis.Client) *RoomStorage {
	return &RoomStorage{rds: rds}
}

// 获取房间名 [ws:sid:room:房间类型:房间号]
func (r *RoomStorage) key(opts *RoomOption) string {
	return fmt.Sprintf("ws:%s:r:%s:%s", opts.Sid, opts.RoomType, opts.Number)
}

// Add 添加房间成员
func (r *RoomStorage) Add(ctx context.Context, opt *RoomOption) error {

	key := r.key(opt)

	err := r.rds.SAdd(ctx, key, opt.Cid).Err()
	if err == nil {
		r.rds.Expire(ctx, key, time.Hour*24*7)
	}

	return err
}

func (r *RoomStorage) BatchAdd(ctx context.Context, opts []*RoomOption) error {

	pipeline := r.rds.Pipeline()
	for _, opt := range opts {
		key := r.key(opt)
		if err := pipeline.SAdd(ctx, key, opt.Cid).Err(); err == nil {
			pipeline.Expire(ctx, key, time.Hour*24*7)
		}
	}

	_, err := pipeline.Exec(ctx)
	return err
}

// Del 删除房间成员
func (r *RoomStorage) Del(ctx context.Context, opts *RoomOption) error {
	return r.rds.SRem(ctx, r.key(opts), opts.Cid).Err()
}

func (r *RoomStorage) BatchDel(ctx context.Context, opts []*RoomOption) error {

	pipeline := r.rds.Pipeline()
	for _, opt := range opts {
		pipeline.SRem(ctx, r.key(opt), opt.Cid)
	}

	_, err := pipeline.Exec(ctx)
	return err
}

// All 获取所有房间成员
func (r *RoomStorage) All(ctx context.Context, opts *RoomOption) []int64 {

	arr := r.rds.SMembers(ctx, r.key(opts)).Val()

	cids := make([]int64, 0, len(arr))
	for _, val := range arr {
		if cid, err := strconv.ParseInt(val, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}
