package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"strconv"
)

type WsClientSession struct {
	rds    *redis.Client
	conf   *config.Config
	server *ServerRunID
}

func NewWsClientSession(
	rds *redis.Client,
	conf *config.Config,
	server *ServerRunID,
) *WsClientSession {
	return &WsClientSession{rds, conf, server}
}

func (w *WsClientSession) getChannelClientKey(sid string, channel string) string {
	return fmt.Sprintf("ws:%s:channel:%s:client", sid, channel)
}

func (w *WsClientSession) getChannelUserKey(sid string, channel string, uid string) string {
	return fmt.Sprintf("ws:%s:channel:%s:user:%s", sid, channel, uid)
}

// Set 设置客户端与用户绑定关系
// channel  渠道分组
// fd       客户端连接ID
// id       用户ID
func (w *WsClientSession) Set(ctx context.Context, channel string, fd string, id int) {
	w.rds.HSet(ctx, w.getChannelClientKey(w.conf.GetSid(), channel), fd, id)

	w.rds.SAdd(ctx, w.getChannelUserKey(w.conf.GetSid(), channel, strconv.Itoa(id)), fd)
}

// Del 删除客户端与用户绑定关系
// channel  渠道分组
// fd     客户端连接ID
func (w *WsClientSession) Del(ctx context.Context, channel string, fd string) {
	KeyName := w.getChannelClientKey(w.conf.GetSid(), channel)

	uid, _ := w.rds.HGet(ctx, KeyName, fd).Result()

	w.rds.HDel(ctx, KeyName, fd)

	w.rds.SRem(ctx, w.getChannelUserKey(w.conf.GetSid(), channel, uid), fd)
}

// IsOnline 判断客户端是否在线[所有部署机器]
// channel  渠道分组
// id       用户ID
func (w *WsClientSession) IsOnline(ctx context.Context, channel string, id string) bool {
	for _, sid := range w.server.GetServerRunIdAll(ctx, 1) {
		val, err := w.rds.SCard(ctx, w.getChannelUserKey(sid, channel, id)).Result()
		if err == nil && val > 0 {
			return true
		}
	}

	return false
}
