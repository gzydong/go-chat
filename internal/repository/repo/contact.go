package repo

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

type Contact struct {
	*Base
	cache    *cache.ContactRemark
	relation *cache.Relation
}

func NewContact(baseDao *Base, cache *cache.ContactRemark, relation *cache.Relation) *Contact {
	return &Contact{Base: baseDao, cache: cache, relation: relation}
}

func (repo *Contact) Remarks(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	if !repo.cache.IsExist(ctx, uid) {
		_ = repo.LoadContactCache(ctx, uid)
	}

	return repo.cache.MGet(ctx, uid, fids)
}

// IsFriend 判断是否为好友关系
func (repo *Contact) IsFriend(ctx context.Context, uid int, friendId int, cache bool) bool {

	if cache && repo.relation.IsContactRelation(ctx, uid, friendId) == nil {
		return true
	}

	sql := `SELECT count(1) from contact where ((user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)) and status = 1`

	var count int
	if err := repo.Db().Raw(sql, uid, friendId, friendId, uid).Scan(&count).Error; err != nil {
		return false
	}

	if count == 2 {
		repo.relation.SetContactRelation(ctx, uid, friendId)
	} else {
		repo.relation.DelContactRelation(ctx, uid, friendId)
	}

	return count == 2
}

func (repo *Contact) GetFriendRemark(ctx context.Context, uid int, friendId int) string {

	if repo.cache.IsExist(ctx, uid) {
		return repo.cache.Get(ctx, uid, friendId)
	}

	info := &model.Contact{}
	repo.db.First(info, "user_id = ? and friend_id = ?", uid, friendId)

	return info.Remark
}

func (repo *Contact) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return repo.cache.Set(ctx, uid, friendId, remark)
}

func (repo *Contact) LoadContactCache(ctx context.Context, uid int) error {

	sql := `SELECT friend_id, remark FROM contact WHERE user_id = ? and status = 1`

	var contacts []*model.Contact
	if err := repo.db.Raw(sql, uid).Scan(&contacts).Error; err != nil {
		return err
	}

	items := make(map[string]interface{})
	for _, value := range contacts {
		if len(value.Remark) > 0 {
			items[fmt.Sprintf("%d", value.FriendId)] = value.Remark
		}
	}

	_ = repo.cache.MSet(ctx, uid, items)

	return nil
}
