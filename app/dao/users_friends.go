package dao

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-chat/app/model"
)

type UsersFriendsDao struct {
	*BaseDao
	rds *redis.Client
}

func NewUsersFriends(base *BaseDao, rds *redis.Client) *UsersFriendsDao {
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
	_, err := dao.BaseUpdate(&model.UsersFriends{}, gin.H{"user_id": uid, "friend_id": friendId}, gin.H{"remark": remark})
	return err
}
