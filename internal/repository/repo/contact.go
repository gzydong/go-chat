package repo

import (
	"context"
	"strconv"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Contact struct {
	ichat.Repo[model.Contact]
	cache    *cache.ContactRemark
	relation *cache.Relation
}

func NewContact(db *gorm.DB, cache *cache.ContactRemark, relation *cache.Relation) *Contact {
	return &Contact{Repo: ichat.NewRepo[model.Contact](db), cache: cache, relation: relation}
}

func (c *Contact) Remarks(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	if !c.cache.Exist(ctx, uid) {
		_ = c.LoadContactCache(ctx, uid)
	}

	return c.cache.MGet(ctx, uid, fids)
}

// IsFriend 判断是否为好友关系
func (c *Contact) IsFriend(ctx context.Context, uid int, friendId int, cache bool) bool {

	if cache && c.relation.IsContactRelation(ctx, uid, friendId) == nil {
		return true
	}

	count, err := c.Repo.QueryCount(ctx, "((user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)) and status = ?", uid, friendId, friendId, uid, model.ContactStatusNormal)
	if err != nil {
		return false
	}

	if count == 2 {
		c.relation.SetContactRelation(ctx, uid, friendId)
	} else {
		c.relation.DelContactRelation(ctx, uid, friendId)
	}

	return count == 2
}

func (c *Contact) GetFriendRemark(ctx context.Context, uid int, friendId int) string {

	if c.cache.Exist(ctx, uid) {
		return c.cache.Get(ctx, uid, friendId)
	}

	var remark string
	c.Repo.Model(ctx).Where("user_id = ? and friend_id = ?", uid, friendId).Pluck("remark", &remark)

	return remark
}

func (c *Contact) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return c.cache.Set(ctx, uid, friendId, remark)
}

func (c *Contact) LoadContactCache(ctx context.Context, uid int) error {

	all, err := c.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Select("friend_id,remark").Where("user_id = ? and status = ?", uid, model.ContactStatusNormal)
	})

	if err != nil {
		return err
	}

	items := make(map[string]any)
	for _, value := range all {
		if len(value.Remark) > 0 {
			items[strconv.Itoa(value.FriendId)] = value.Remark
		}
	}

	return c.cache.MSet(ctx, uid, items)
}
