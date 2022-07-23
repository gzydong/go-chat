package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// ContactRemark 联系人备注缓存
type ContactRemark struct {
	rds *redis.Client
}

func NewContactRemark(rds *redis.Client) *ContactRemark {
	return &ContactRemark{rds: rds}
}

func (c *ContactRemark) name(uid int) string {
	return fmt.Sprintf("contact:remark:uid_%d", uid)
}

func (c *ContactRemark) Get(ctx context.Context, uid int, fid int) string {
	return c.rds.HGet(ctx, c.name(uid), fmt.Sprintf("%d", fid)).Val()
}

func (c *ContactRemark) MGet(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	values := make([]string, 0, len(fids))
	for _, value := range fids {
		values = append(values, fmt.Sprintf("%d", value))
	}

	items := c.rds.HMGet(ctx, c.name(uid), values...).Val()

	remarks := make(map[int]string)
	for k, v := range fids {
		if items[k] != nil {
			remarks[v] = items[k].(string)
		}
	}

	return remarks, nil
}

// Set 设置备注
func (c *ContactRemark) Set(ctx context.Context, uid int, friendId int, value string) error {

	if c.rds.Exists(ctx, c.name(uid)).Val() == 1 {
		return c.rds.HSet(ctx, c.name(uid), friendId, value).Err()
	}

	return nil
}

func (c *ContactRemark) IsExist(ctx context.Context, uid int) bool {
	return c.rds.Exists(ctx, c.name(uid)).Val() == 1
}

// MSet 批量设置备注
func (c *ContactRemark) MSet(ctx context.Context, uid int, values map[string]interface{}) error {

	c.rds.HSet(ctx, c.name(uid), values)
	c.rds.Expire(ctx, c.name(uid), time.Hour*24)

	return nil
}
