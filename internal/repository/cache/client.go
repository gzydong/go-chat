package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go-chat/config"
)

type ClientStorage struct {
	redis   *redis.Client
	config  *config.Config
	storage *ServerStorage
}

func NewClientStorage(redis *redis.Client, config *config.Config, storage *ServerStorage) *ClientStorage {
	return &ClientStorage{redis: redis, config: config, storage: storage}
}

// Set 设置客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
// @params id       用户ID
func (w *ClientStorage) Set(ctx context.Context, channel string, fd string, uid int) error {
	sid := w.config.ServerId()
	_, err := w.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, w.clientKey(sid, channel), fd, uid)
		pipe.SAdd(ctx, w.userKey(sid, channel, strconv.Itoa(uid)), fd)
		return nil
	})
	return err
}

// Del 删除客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
func (w *ClientStorage) Del(ctx context.Context, channel, fd string) error {
	sid := w.config.ServerId()
	key := w.clientKey(sid, channel)
	uid, _ := w.redis.HGet(ctx, key, fd).Result()
	_, err := w.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HDel(ctx, key, fd)
		pipe.SRem(ctx, w.userKey(sid, channel, uid), fd)
		return nil
	})
	return err
}

// IsOnline 判断客户端是否在线[所有部署机器]
// @params channel  渠道分组
// @params uid      用户ID
func (w *ClientStorage) IsOnline(ctx context.Context, channel, uid string) bool {
	for _, sid := range w.storage.All(ctx, 1) {
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
	val, err := w.redis.SCard(ctx, w.userKey(sid, channel, uid)).Result()

	return err == nil && val > 0
}

// GetUidFromClientIds 获取当前节点用户ID关联的客户端ID
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (w *ClientStorage) GetUidFromClientIds(ctx context.Context, sid, channel, uid string) []int64 {
	cids := make([]int64, 0)

	items, err := w.redis.SMembers(ctx, w.userKey(sid, channel, uid)).Result()
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
	uid, err := w.redis.HGet(ctx, w.clientKey(sid, channel), cid).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(uid, 10, 64)
}

func (w *ClientStorage) Bind(ctx context.Context, channel string, clientId int64, uid int) error {
	return w.Set(ctx, channel, strconv.FormatInt(clientId, 10), uid)
}

func (w *ClientStorage) UnBind(ctx context.Context, channel string, clientId int64) error {
	return w.Del(ctx, channel, strconv.FormatInt(clientId, 10))
}

func (w *ClientStorage) clientKey(sid, channel string) string {
	return fmt.Sprintf("ws:%s:channel:%s:redis", sid, channel)
}

func (w *ClientStorage) userKey(sid, channel, uid string) string {
	return fmt.Sprintf("ws:%s:channel:%s:user:%s", sid, channel, uid)
}
