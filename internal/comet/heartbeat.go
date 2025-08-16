package comet

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/longnet"
	"go-chat/internal/repository/cache"
)

type Heartbeat struct {
	ServerStorage *cache.ServerStorage
	Redis         *redis.Client
}

type ServerInfo struct {
	StartAt         string `json:"start_at"`           // 启动时间
	ActiveAt        string `json:"active_at"`          // 心跳时间
	CurrConnNum     int32  `json:"curr_conn_num"`      // 当前连接数
	CurrConnUserNum int32  `json:"curr_conn_user_num"` // 当前用户数
	AutoConnId      int64  `json:"auto_conn_id"`       // 自动连接ID
	SendMessageNum  int64  `json:"send_message_num"`   // 发送消息数
	RecvMessageNum  int64  `json:"recv_message_num"`   // 接收消息数
}

func (h *Heartbeat) Start(ctx context.Context, serv longnet.IServer) error {
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	h.save(ctx, serv)

	for {
		select {
		case <-ctx.Done():
			h.Redis.HDel(context.Background(), "im:server_infos", serv.ServerId())
			return nil
		case <-timer.C:
			h.save(ctx, serv)
		}
	}
}

func (h *Heartbeat) save(ctx context.Context, serv longnet.IServer) {
	info := &ServerInfo{
		StartAt: time.Now().Format(time.DateTime),
	}

	value := h.Redis.HGet(ctx, "im:server_infos", serv.ServerId()).Val()
	if value != "" {
		_ = jsonutil.Unmarshal(value, info)
	}

	info.ActiveAt = time.Now().Format(time.DateTime)
	info.CurrConnNum = serv.SessionManager().GetSessionNum()
	info.CurrConnUserNum = serv.SessionManager().GetSessionUserNum()
	info.SendMessageNum = 0
	info.RecvMessageNum = 0

	h.Redis.HSet(ctx, "im:server_infos", serv.ServerId(), jsonutil.Encode(info))
}
