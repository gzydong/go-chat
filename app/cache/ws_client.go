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
func (w *WsClient) Set(channel string, uuid string, id int) {
	flag := fmt.Sprintf("ws:channel:%s:client", channel)
	Rdb.HSet(flag, uuid, id)

	flag = fmt.Sprintf("ws:channel:%s:user:%d", channel, id)
	Rdb.SAdd(flag, uuid)
}

// Del 删除客户端与用户绑定关系
func (w *WsClient) Del(channel string, uuid string) {
	flag := fmt.Sprintf("ws:channel:%s", channel)

	id, _ := Rdb.HGet(flag, uuid).Result()

	Rdb.HDel(flag, uuid)

	flag = fmt.Sprintf("ws:channel:%s:user:%s", channel, id)
	Rdb.SRem(flag, uuid)
}
