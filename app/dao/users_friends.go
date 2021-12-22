package dao

import (
	"context"
	"fmt"
	"go-chat/app/model"
	"time"
)

type UsersFriendsDao struct {
	*BaseDao
}

func NewUsersFriends(base *BaseDao) *UsersFriendsDao {
	return &UsersFriendsDao{base}
}

// IsFriend 判断是否为好友关系
func (dao *UsersFriendsDao) IsFriend(ctx context.Context, uid int, friendId int) bool {
	sql := `SELECT count(1) from users_friends where ((user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)) and status = 1`

	var count int
	if err := dao.Db().Raw(sql, uid, friendId, friendId, uid).Scan(&count).Error; err != nil {
		return false
	}

	return count == 2
}

func (dao *UsersFriendsDao) GetFriendRemark(ctx context.Context, uid int, friendId int, isCache bool) string {

	if isCache {
		remark := dao.rds.HGet(ctx, fmt.Sprintf("rds:hash:friend-remark:%d", uid), fmt.Sprintf("%d_%d", uid, friendId)).Val()
		if remark != "" {
			return remark
		}
	}

	remark := ""
	err := dao.Db().Model(&model.UsersFriends{}).Select("remark").Where("user_id = ? and friend_id = ?", uid, friendId).Scan(&remark).Error
	if err != nil {
		_ = dao.SetFriendRemark(ctx, uid, friendId, remark)
	}

	return remark
}

func (dao *UsersFriendsDao) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	err := dao.rds.HSet(ctx, fmt.Sprintf("rds:hash:friend-remark:%d", uid), fmt.Sprintf("%d_%d", uid, friendId), remark).Err()
	if err == nil {
		dao.rds.Expire(ctx, fmt.Sprintf("rds:hash:friend-remark:%d", uid), 72*time.Hour)
	}

	return err
}
