package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/entity"
)

type Room struct {
	rds *redis.Client
}

type RoomOption struct {
	Channel  string          // 渠道分类
	RoomType entity.RoomType // 房间类型
	Number   string          // 房间号
	Sid      string          // 网关ID
	Cid      int64           // 客户端ID
}

func NewRoom(rds *redis.Client) *Room {
	return &Room{rds: rds}
}

// 获取房间名 [ws:sid:room:房间类型:房间号]
func (room *Room) key(opts *RoomOption) string {
	return fmt.Sprintf("ws:%s:room:%s:%s", opts.Sid, opts.RoomType, opts.Number)
}

// Add 添加房间成员
func (room *Room) Add(ctx context.Context, opts *RoomOption) error {

	key := room.key(opts)

	err := room.rds.SAdd(ctx, key, opts.Cid).Err()
	if err == nil {
		room.rds.Expire(ctx, key, time.Hour*24*7)
	}

	return err
}

// Del 删除房间成员
func (room *Room) Del(ctx context.Context, opts *RoomOption) error {
	return room.rds.SRem(ctx, room.key(opts), opts.Cid).Err()
}

// All 获取所有房间成员
func (room *Room) All(ctx context.Context, opts *RoomOption) []int64 {

	arr := room.rds.SMembers(ctx, room.key(opts)).Val()

	cids := make([]int64, 0, len(arr))
	for _, val := range arr {
		if cid, err := strconv.ParseInt(val, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}
