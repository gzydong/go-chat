package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go-chat/config"
)

type ClientStorage struct {
	redis  *redis.Client
	config *config.Config
	server *SidServer
}

func NewClientStorage(redis *redis.Client, config *config.Config, server *SidServer) *ClientStorage {
	return &ClientStorage{redis: redis, config: config, server: server}
}

func (w *ClientStorage) getClientKey(sid, channel string) string {
	return fmt.Sprintf("ws:%s:channel:%s:client", sid, channel)
}

func (w *ClientStorage) getUserKey(sid, channel, uid string) string {
	return fmt.Sprintf("ws:%s:channel:%s:user:%s", sid, channel, uid)
}

// Set 设置客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
// @params id       用户ID
func (w *ClientStorage) Set(ctx context.Context, channel string, fd string, uid int) {
	w.redis.HSet(ctx, w.getClientKey(w.config.ServerId(), channel), fd, uid)

	w.redis.SAdd(ctx, w.getUserKey(w.config.ServerId(), channel, strconv.Itoa(uid)), fd)
}

// Del 删除客户端与用户绑定关系
// @params channel  渠道分组
// @params fd     客户端连接ID
func (w *ClientStorage) Del(ctx context.Context, channel, fd string) {
	KeyName := w.getClientKey(w.config.ServerId(), channel)

	uid, _ := w.redis.HGet(ctx, KeyName, fd).Result()

	w.redis.HDel(ctx, KeyName, fd)

	w.redis.SRem(ctx, w.getUserKey(w.config.ServerId(), channel, uid), fd)
}

// IsOnline 判断客户端是否在线[所有部署机器]
// @params channel  渠道分组
// @params uid      用户ID
func (w *ClientStorage) IsOnline(ctx context.Context, channel, uid string) bool {
	for _, sid := range w.server.All(ctx, 1) {
		if w.IsCurrentServerOnline(ctx, sid, channel, uid) {
			return true
		}
	}

	return false
}

// IsCurrentServerOnline 判断当前节点是否在线
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (w *ClientStorage) IsCurrentServerOnline(ctx context.Context, sid, channel, uid string) bool {
	val, err := w.redis.SCard(ctx, w.getUserKey(sid, channel, uid)).Result()

	return err == nil && val > 0
}

// GetUidFromClientIds 获取当前节点用户ID关联的客户端ID
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (w *ClientStorage) GetUidFromClientIds(ctx context.Context, sid, channel, uid string) []int64 {
	cids := make([]int64, 0)

	items, err := w.redis.SMembers(ctx, w.getUserKey(sid, channel, uid)).Result()
	if err != nil {
		return cids
	}

	for _, cid := range items {
		if cid, err := strconv.ParseInt(cid, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}

// GetClientIdFromUid 获取客户端ID关联的用户ID
// @params sid     服务节点ID
// @params channel 渠道分组
// @params cid     客户端ID
func (w *ClientStorage) GetClientIdFromUid(ctx context.Context, sid, channel, cid string) (int64, error) {
	uid, err := w.redis.HGet(ctx, w.getClientKey(sid, channel), cid).Result()
	if err != nil {
		return 0, err
	}

	if value, err := strconv.ParseInt(uid, 10, 64); err != nil {
		return value, nil
	} else {
		return 0, err
	}
}

func (w *ClientStorage) Bind(ctx context.Context, channel string, clientId int64, uid int) {
	w.Set(ctx, channel, fmt.Sprintf("%d", clientId), uid)
}

func (w *ClientStorage) UnBind(ctx context.Context, channel string, clientId int64) {
	w.Del(ctx, channel, fmt.Sprintf("%d", clientId))
}
