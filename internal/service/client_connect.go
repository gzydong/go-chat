package service

import (
	"context"
	"strconv"

	"go-chat/internal/repository/cache"
)

var _ IClientConnectService = (*ClientConnectService)(nil)

type IClientConnectService interface {
	// Bind 绑定用户和客户端
	Bind(ctx context.Context, sid, channel string, clientId int64, uid int) error
	// UnBind 解除用户和客户端的绑定
	UnBind(ctx context.Context, sid, channel string, clientId int64) error
	// IsUidOnline 检查用户是否在线(不区分服务ID)
	IsUidOnline(ctx context.Context, channel string, uid int) (bool, error)
	// IsUidOnlineBySid 检查用户是否在线(指定服务ID)
	IsUidOnlineBySid(ctx context.Context, sid, channel string, uid int) (bool, error)
	// GetUidByClientId 获取客户端绑定的用户
	GetUidByClientId(ctx context.Context, sid, channel string, clientId int64) (int64, error)
	// GetUidFromClientIds 获取用户绑定的客户端
	GetUidFromClientIds(ctx context.Context, sid, channel string, uid int) ([]int64, error)
}

// ClientConnectService 客户端连接管理服务
type ClientConnectService struct {
	Storage *cache.ClientStorage
}

func (c *ClientConnectService) Bind(ctx context.Context, sid, channel string, clientId int64, uid int) error {
	return c.Storage.Bind(ctx, sid, channel, clientId, uid)
}

func (c *ClientConnectService) UnBind(ctx context.Context, sid, channel string, clientId int64) error {
	return c.Storage.UnBind(ctx, sid, channel, clientId)
}

func (c *ClientConnectService) IsUidOnline(ctx context.Context, channel string, uid int) (bool, error) {
	return c.Storage.IsOnline(ctx, channel, strconv.Itoa(uid)), nil
}

func (c *ClientConnectService) IsUidOnlineBySid(ctx context.Context, sid, channel string, uid int) (bool, error) {
	return c.Storage.IsCurrentServerOnline(ctx, sid, channel, strconv.Itoa(uid)), nil
}

func (c *ClientConnectService) GetUidByClientId(ctx context.Context, sid, channel string, clientId int64) (int64, error) {
	return c.Storage.GetClientIdFromUid(ctx, sid, channel, clientId)
}

func (c *ClientConnectService) GetUidFromClientIds(ctx context.Context, sid, channel string, uid int) ([]int64, error) {
	return c.Storage.GetUidFromClientIds(ctx, sid, channel, strconv.Itoa(uid)), nil
}
