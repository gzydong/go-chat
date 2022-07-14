package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type UnreadStorage struct {
	rds *redis.Client
}

func NewUnreadStorage(rds *redis.Client) *UnreadStorage {
	return &UnreadStorage{rds}
}

func (u *UnreadStorage) name(receive int) string {
	return fmt.Sprintf("talk:unread_msg:uid_%d", receive)
}

// Increment 消息未读数自增
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u *UnreadStorage) Increment(ctx context.Context, mode, sender, receive int) {
	u.rds.HIncrBy(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender), 1)
}

// Get 获取消息未读数
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u *UnreadStorage) Get(ctx context.Context, mode, sender, receive int) int {
	val, _ := u.rds.HGet(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender)).Int()

	return val
}

func (u UnreadStorage) Del(ctx context.Context, mode, sender, receive int) {
	u.rds.HDel(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender))
}

// Reset 消息未读数重置
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u *UnreadStorage) Reset(ctx context.Context, mode, sender, receive int) {
	u.rds.HSet(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender), 0)
}

func (u *UnreadStorage) GetAll(ctx context.Context, receive int) map[string]int {
	items := make(map[string]int)
	for k, v := range u.rds.HGetAll(ctx, u.name(receive)).Val() {
		items[k], _ = strconv.Atoi(v)
	}

	return items
}
