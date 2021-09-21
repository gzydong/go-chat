package cache

import (
	"fmt"
)

type WsClient struct {
	CacheKey string
}

func NewWsClient() *WsClient {
	return &WsClient{}
}

// Set 设置客户端与用户绑定关系
// channel  渠道分组
// uuid     客户端连接ID
// id       用户ID
func (w *WsClient) Set(channel string, uuid string, id int) {
	flag := fmt.Sprintf("ws:channel:%s:client", channel)
	Rdb.HSet(flag, uuid, id)

	flag = fmt.Sprintf("ws:channel:%s:user:%d", channel, id)
	Rdb.SAdd(flag, uuid)
}

// Del 删除客户端与用户绑定关系
// channel  渠道分组
// uuid     客户端连接ID
func (w *WsClient) Del(channel string, uuid string) {
	flag := fmt.Sprintf("ws:channel:%s", channel)

	id, _ := Rdb.HGet(flag, uuid).Result()

	Rdb.HDel(flag, uuid)

	flag = fmt.Sprintf("ws:channel:%s:user:%s", channel, id)
	Rdb.SRem(flag, uuid)
}

// IsOnline 判断客户端是否在线[当前机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnline(channel string, id string) bool {
	flag := fmt.Sprintf("ws:channel:%s:user:%s", channel, id)

	val, err := Rdb.SCard(flag).Result()
	if err != nil {
		return false
	}

	return val > 0
}

// IsOnlineAll 判断客户端是否在线[所有部署机器]
// channel  渠道分组
// id       用户ID
func (w *WsClient) IsOnlineAll(channel string, id string) bool {
	flag := fmt.Sprintf("ws:channel:%s:user:%s", channel, id)

	val, err := Rdb.SCard(flag).Result()
	if err != nil {
		return false
	}

	return val > 0
}
