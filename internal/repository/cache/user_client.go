package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/jsonutil"
)

type UserClient struct {
	redis *redis.Client
}

// hash im:user_clients:#user_id
//	  serv_id1:1 => {"activity_at":1750304874}
//	  serv_id1:2 => {"activity_at":1750304874}
//	  serv_id2:3 => {"activity_at":1750304874}

func NewUserClient(rds *redis.Client) *UserClient {
	return &UserClient{
		redis: rds,
	}
}

func (u *UserClient) key(uid int64) string {
	return fmt.Sprintf("im:user_clients:%d", uid)
}

func (u *UserClient) Bind(ctx context.Context, serverId string, clientId int64, uid int64) error {
	pipeline := u.redis.Pipeline()

	pipeline.HMSet(ctx, u.key(uid), fmt.Sprintf("%s:%d", serverId, clientId), jsonutil.Marshal(map[string]int64{
		"active_at": time.Now().Unix(),
	}))

	pipeline.Expire(ctx, u.key(uid), time.Hour*24*3)
	_, err := pipeline.Exec(ctx)
	return err
}

func (u *UserClient) UnBind(ctx context.Context, serverId string, clientId int64, uid int64) error {
	return u.redis.HDel(ctx, u.key(uid), fmt.Sprintf("%s:%d", serverId, clientId)).Err()
}

type Client struct {
	ServerId string `json:"server_id"`
	ClientId int64  `json:"client_id"`
	ActiveAt int64  `json:"active_at"`
}

func (u *UserClient) GetClientList(ctx context.Context, uid int64) ([]*Client, error) {
	items := make([]*Client, 0)

	resp, err := u.redis.HGetAll(ctx, u.key(uid)).Result()
	if err != nil {
		return nil, err
	}

	for k, v := range resp {
		client := &Client{}
		if err := jsonutil.Unmarshal([]byte(v), client); err != nil {
			continue
		}

		i := strings.Index(k, ":")

		client.ServerId = k[:i]

		cid, _ := strconv.Atoi(k[i+1:])
		client.ClientId = int64(cid)
		items = append(items, client)
	}

	return items, nil
}

func (u *UserClient) IsOnline(ctx context.Context, uid int64) bool {
	clients, err := u.GetClientList(ctx, uid)
	if err != nil {
		return false
	}

	delClients := make([]*Client, 0)
	for _, client := range clients {
		// 超过5分钟则视为离线
		if time.Now().Unix()-client.ActiveAt > 60*5 {
			delClients = append(delClients, client)
		}
	}

	if len(delClients) > 0 {
		pipeline := u.redis.Pipeline()
		for _, client := range delClients {
			pipeline.HDel(ctx, u.key(uid), fmt.Sprintf("%s:%d", client.ServerId, client.ClientId))
		}
		_, _ = pipeline.Exec(ctx)
	}

	for _, client := range clients {
		if client.ActiveAt > time.Now().Unix()-60 {
			return true
		}
	}

	return false
}
