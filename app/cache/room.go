package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/entity"
	"strconv"
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

// 获取房间名 [ws:网关ID:room:房间类型:房间号]
func (room *Room) key(opts *RoomOption) string {
	return fmt.Sprintf("ws:%s:room:%s:%s", opts.Sid, opts.RoomType, opts.Number)
}

// Add 添加房间成员
func (room *Room) Add(ctx context.Context, opts *RoomOption) error {
	return room.rds.SAdd(ctx, room.key(opts), opts.Cid).Err()
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
