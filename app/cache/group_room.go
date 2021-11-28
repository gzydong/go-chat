package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type GroupRoom struct {
	rds *redis.Client
}

func NewGroupRoom(rds *redis.Client) *GroupRoom {
	return &GroupRoom{rds: rds}
}

func (room GroupRoom) key(sid, name string) string {
	return fmt.Sprintf("ws:%s:room:group-room:%s", sid, name)
}

// Add 添加房间成员
func (room *GroupRoom) Add(ctx context.Context, sid string, name string, cid int64) error {
	return room.rds.SAdd(ctx, room.key(sid, name), cid).Err()
}

// Del 删除房间成员
func (room *GroupRoom) Del(ctx context.Context, sid string, name string, cid int64) error {
	return room.rds.SRem(ctx, room.key(sid, name), cid).Err()
}

// All 获取所有房间成员
func (room *GroupRoom) All(ctx context.Context, sid string, name string) []int64 {

	arr := room.rds.SMembers(ctx, room.key(sid, name)).Val()

	cids := make([]int64, 0, len(arr))
	for _, val := range arr {
		if cid, err := strconv.ParseInt(val, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}
