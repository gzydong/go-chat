package dao

import (
	"context"
	"fmt"

	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
)

type IContactDao interface {
	IBaseDao
	IsFriend(ctx context.Context, uid int, friendId int, cache bool) bool
	GetFriendRemark(ctx context.Context, uid int, friendId int) string
	SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error
	Remarks(ctx context.Context, uid int, fids []int) (map[int]string, error)
}

type ContactDao struct {
	*BaseDao
	cache    *cache.ContactRemark
	relation *cache.Relation
}

func NewContactDao(baseDao *BaseDao, cache *cache.ContactRemark, relation *cache.Relation) *ContactDao {
	return &ContactDao{BaseDao: baseDao, cache: cache, relation: relation}
}

func (dao *ContactDao) Remarks(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	if !dao.cache.IsExist(ctx, uid) {
		_ = dao.LoadContactCache(ctx, uid)
	}

	return dao.cache.MGet(ctx, uid, fids)
}

// IsFriend 判断是否为好友关系
func (dao *ContactDao) IsFriend(ctx context.Context, uid int, friendId int, cache bool) bool {

	if cache && dao.relation.IsContactRelation(ctx, uid, friendId) == nil {
		return true
	}

	sql := `SELECT count(1) from contact where ((user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)) and status = 1`

	var count int
	if err := dao.Db().Raw(sql, uid, friendId, friendId, uid).Scan(&count).Error; err != nil {
		return false
	}

	if count == 2 {
		dao.relation.SetContactRelation(ctx, uid, friendId)
	} else {
		dao.relation.DelContactRelation(ctx, uid, friendId)
	}

	return count == 2
}

func (dao *ContactDao) GetFriendRemark(ctx context.Context, uid int, friendId int) string {

	if dao.cache.IsExist(ctx, uid) {
		return dao.cache.Get(ctx, uid, friendId)
	}

	info := &model.Contact{}
	dao.db.First(info, "user_id = ? and friend_id = ?", uid, friendId)

	return info.Remark
}

func (dao *ContactDao) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return dao.cache.Set(ctx, uid, friendId, remark)
}

func (dao *ContactDao) LoadContactCache(ctx context.Context, uid int) error {

	sql := `SELECT friend_id, remark FROM contact WHERE user_id = ? and status = 1`

	var contacts []*model.Contact
	if err := dao.db.Raw(sql, uid).Scan(&contacts).Error; err != nil {
		return err
	}

	items := make(map[string]interface{})
	for _, value := range contacts {
		if len(value.Remark) > 0 {
			items[fmt.Sprintf("%d", value.FriendId)] = value.Remark
		}
	}

	_ = dao.cache.MSet(ctx, uid, items)

	return nil
}
