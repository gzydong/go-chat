package repo

import (
	"context"
	"fmt"

	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
)

type IContact interface {
	IBase
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

func NewContact(base *Base, cache *cache.ContactRemark, relation *cache.Relation) *Contact {
	return &Contact{Base: base, cache: cache, relation: relation}
}

func (c *Contact) Remarks(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	if !c.cache.IsExist(ctx, uid) {
		_ = c.LoadContactCache(ctx, uid)
	}

	return c.cache.MGet(ctx, uid, fids)
}

// IsFriend 判断是否为好友关系
func (c *Contact) IsFriend(ctx context.Context, uid int, friendId int, cache bool) bool {

	if cache && c.relation.IsContactRelation(ctx, uid, friendId) == nil {
		return true
	}

	sql := `SELECT count(1) from contact where ((user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)) and status = 1`

	var count int
	if err := c.Db.WithContext(ctx).Raw(sql, uid, friendId, friendId, uid).Scan(&count).Error; err != nil {
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

	if c.cache.IsExist(ctx, uid) {
		return c.cache.Get(ctx, uid, friendId)
	}

	info := &model.Contact{}
	c.Db.WithContext(ctx).First(info, "user_id = ? and friend_id = ?", uid, friendId)

	return info.Remark
}

func (c *Contact) SetFriendRemark(ctx context.Context, uid int, friendId int, remark string) error {
	return c.cache.Set(ctx, uid, friendId, remark)
}

func (c *Contact) LoadContactCache(ctx context.Context, uid int) error {

	sql := `SELECT friend_id, remark FROM contact WHERE user_id = ? and status = 1`

	var contacts []*model.Contact
	if err := c.Db.WithContext(ctx).Raw(sql, uid).Scan(&contacts).Error; err != nil {
		return err
	}

	items := make(map[string]interface{})
	for _, value := range contacts {
		if len(value.Remark) > 0 {
			items[fmt.Sprintf("%d", value.FriendId)] = value.Remark
		}
	}

	_ = c.cache.MSet(ctx, uid, items)

	return nil
}