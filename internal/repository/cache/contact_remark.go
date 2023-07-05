package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ContactRemark 联系人备注缓存
type ContactRemark struct {
	redis *redis.Client
}

func NewContactRemark(redis *redis.Client) *ContactRemark {
	return &ContactRemark{redis: redis}
}

func (c *ContactRemark) Get(ctx context.Context, uid int, fid int) string {
	return c.redis.HGet(ctx, c.name(uid), fmt.Sprintf("%d", fid)).Val()
}

func (c *ContactRemark) MGet(ctx context.Context, uid int, fids []int) (map[int]string, error) {

	values := make([]string, 0, len(fids))
	for _, value := range fids {
		values = append(values, strconv.Itoa(value))
	}

	remarks := make(map[int]string)

	items, err := c.redis.HMGet(ctx, c.name(uid), values...).Result()
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return remarks, nil
	}

	for k, v := range fids {
		if items[k] != nil {
			remarks[v] = items[k].(string)
		}
	}

	return remarks, nil
}

func (c *ContactRemark) Set(ctx context.Context, uid int, friendId int, value string) error {

	if c.Exist(ctx, uid) {
		return c.redis.HSet(ctx, c.name(uid), friendId, value).Err()
	}

	return nil
}

func (c *ContactRemark) MSet(ctx context.Context, uid int, values map[string]any) error {
	_, err := c.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, c.name(uid), values)
		pipe.Expire(ctx, c.name(uid), 12*time.Hour)
		return nil
	})
	return err
}

func (c *ContactRemark) Exist(ctx context.Context, uid int) bool {
	return c.redis.Exists(ctx, c.name(uid)).Val() == 1
}

func (c *ContactRemark) name(uid int) string {
	return fmt.Sprintf("im:contact:remark:uid_%d", uid)
}
