package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"strconv"
)

type WsClient struct {
	Redis  *redis.Client
	Conf   *config.Config
	Server *ServerRunID
}

func (w *WsClient) getChannelClientKey(sid string, channel string) string {
	return fmt.Sprintf("ws:%s:channel:%s:client", sid, channel)
}

func (w *WsClient) getChannelUserKey(sid string, channel string, uid string) string {
	return fmt.Sprintf("ws:%s:channel:%s:user:%s", sid, channel, uid)
}

// Set 设置客户端与用户绑定关系
// channel  渠道分组
// fd       客户端连接ID
// id       用户ID
func (w *WsClient) Set(ctx context.Context, channel string, fd string, id int) {
	w.Redis.HSet(ctx, w.getChannelClientKey(w.Conf.GetSid(), channel), fd, id)

	w.Redis.SAdd(ctx, w.getChannelUserKey(w.Conf.GetSid(), channel, strconv.Itoa(id)), fd)
}

// Del 删除客户端与用户绑定关系
// channel  渠道分组
// fd     客户端连接ID
func (w *WsClient) Del(ctx context.Context, channel string, fd string) {
	KeyName := w.getChannelClientKey(w.Conf.GetSid(), channel)

	userId, _ := w.Redis.HGet(ctx, KeyName, fd).Result()

	w.Redis.HDel(ctx, KeyName, fd)

	w.Redis.SRem(ctx, w.getChannelUserKey(w.Conf.GetSid(), channel, userId), fd)
}

// IsOnline 判断客户端是否在线[当前机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnline(ctx context.Context, channel string, id string) bool {
	val, err := w.Redis.SCard(ctx, w.getChannelUserKey(w.Conf.GetSid(), channel, id)).Result()

	return err == nil && val > 0
}

// IsOnlineAll 判断客户端是否在线[所有部署机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnlineAll(ctx context.Context, channel string, id string) bool {
	for _, sid := range w.Server.GetServerRunIdAll(ctx, 1) {
		key := w.getChannelUserKey(sid, channel, id)
		val, err := w.Redis.SCard(ctx, key).Result()

		if err == nil && val > 0 {
			return true
		}
	}

	return false
}
