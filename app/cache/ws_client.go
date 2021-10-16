package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"strconv"
)

type WsClient struct {
	Redis *redis.Client
	Conf  *config.Config
}

func (w *WsClient) getChannelClientKey(channel string) string {
	return fmt.Sprintf("ws:%s:channel:%s:client", w.Conf.Server.ServerId, channel)
}

func (w *WsClient) getChannelUserKey(channel string, uid string) string {
	return fmt.Sprintf("ws:%s:channel:%s:user:%s", w.Conf.Server.ServerId, channel, uid)
}

// Set 设置客户端与用户绑定关系
// channel  渠道分组
// fd       客户端连接ID
// id       用户ID
func (w *WsClient) Set(ctx context.Context, channel string, fd string, id int) {
	w.Redis.HSet(ctx, w.getChannelClientKey(channel), fd, id)

	w.Redis.SAdd(ctx, w.getChannelUserKey(channel, strconv.Itoa(id)), fd)
}

// Del 删除客户端与用户绑定关系
// channel  渠道分组
// fd     客户端连接ID
func (w *WsClient) Del(ctx context.Context, channel string, fd string) {
	KeyName := w.getChannelClientKey(channel)

	userId, _ := w.Redis.HGet(ctx, KeyName, fd).Result()

	w.Redis.HDel(ctx, KeyName, fd)

	w.Redis.SRem(ctx, w.getChannelUserKey(channel, userId), fd)
}

// IsOnline 判断客户端是否在线[当前机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnline(ctx context.Context, channel string, id string) bool {
	val, err := w.Redis.SCard(ctx, w.getChannelUserKey(channel, id)).Result()

	return err != nil && val > 0
}

// IsOnlineAll 判断客户端是否在线[所有部署机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnlineAll(ctx context.Context, channel string, id string) bool {
	val, err := w.Redis.SCard(ctx, w.getChannelUserKey(channel, id)).Result()

	return err != nil && val > 0
}
