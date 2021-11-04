package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type UsersFriendsDao struct {
	*Base
	rds *redis.Client
}

func NewUsersFriends(base *Base, rds *redis.Client) *UsersFriendsDao {
	return &UsersFriendsDao{base, rds}
}

func (dao *UsersFriendsDao) GetFriendRemark(ctx context.Context, uid int, friendId int) string {
	res, err := dao.rds.HGet(ctx, "rds:hash:friend-remark", fmt.Sprintf("%d_%d", uid, friendId)).Result()
	if err != nil {
		return res
	}

	return ""
}

func (dao *UsersFriendsDao) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return nil
}
