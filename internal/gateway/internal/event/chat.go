package event

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/gateway/internal/event/chat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/service"
)

type ChatEvent struct {
	redis         *redis.Client
	config        *config.Config
	roomStorage   *cache.RoomStorage
	memberService *service.GroupMemberService
	handler       *chat.Handler
}

func NewChatEvent(redis *redis.Client, config *config.Config, roomStorage *cache.RoomStorage, memberService *service.GroupMemberService, handler *chat.Handler) *ChatEvent {
	return &ChatEvent{redis: redis, config: config, roomStorage: roomStorage, memberService: memberService, handler: handler}
}

// OnOpen 连接成功回调事件
func (d *ChatEvent) OnOpen(client im.IClient) {

	ctx := context.Background()

	// 1.查询用户群列表
	ids := d.memberService.Dao().GetUserGroupIds(context.Background(), client.Uid())

	// 2.客户端加入群房间
	for _, id := range ids {
		_ = d.roomStorage.Add(ctx, &cache.RoomOption{
			Channel:  im.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      d.config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	// 推送上线消息
	d.redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.Uid(),
			"status":  1,
		}),
	}))
}

// OnMessage 消息回调事件
func (d *ChatEvent) OnMessage(client im.IClient, message []byte) {

	// 获取事件名
	event := gjson.GetBytes(message, "event").String()
	if event != "" {
		// 触发事件
		d.handler.Call(context.Background(), client, event, message)
	}
}

// OnClose 连接关闭回调事件
func (d *ChatEvent) OnClose(client im.IClient, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := d.memberService.Dao().GetUserGroupIds(context.Background(), client.Uid())

	// 3.客户端退出群房间
	for _, id := range ids {
		_ = d.roomStorage.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      d.config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	// 推送下线消息
	d.redis.Publish(context.Background(), entity.ImTopicChat, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.Uid(),
			"status":  0,
		}),
	}))
}
