package event

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/gjson"
	"go-chat/internal/business"
	"go-chat/internal/comet/handler/event/chat"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type ChatEvent struct {
	Redis           *redis.Client
	GroupMemberRepo *repo.GroupMember
	MemberService   service.IGroupMemberService
	Handler         *chat.Handler
	RoomStorage     *socket.RoomStorage
	PushMessage     *business.PushMessage
}

// OnOpen 连接成功回调事件
func (c *ChatEvent) OnOpen(client socket.IClient) {
	ctx := context.TODO()

	now := time.Now()

	// 客户端加入群房间
	for _, groupId := range c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid()) {
		_ = c.RoomStorage.Insert(int32(groupId), client.Cid(), now.Unix())
	}

	// 推送上线消息
	_ = c.PushMessage.Push(ctx, entity.ImTopicChat, &entity.SubscribeMessage{
		Event: entity.SubEventContactStatus,
		Payload: jsonutil.Encode(entity.SubEventContactStatusPayload{
			Status: 1,
			UserId: client.Uid(),
		}),
	})
}

// OnMessage 消息回调事件
func (c *ChatEvent) OnMessage(client socket.IClient, message []byte) {
	res := gjson.GetBytes(message, "event")
	if !res.Exists() {
		return
	}

	// 获取事件名
	event := res.String()
	if event != "" {
		// 触发事件
		c.Handler.Call(context.TODO(), client, event, message)
	}
}

// OnClose 连接关闭回调事件
func (c *ChatEvent) OnClose(client socket.IClient, code int, text string) {
	ctx := context.TODO()

	now := time.Now()

	// 客户端退出群房间
	for _, groupId := range c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid()) {
		_ = c.RoomStorage.Delete(int32(groupId), client.Cid(), now.Unix())
	}

	// 推送下线消息
	_ = c.PushMessage.Push(ctx, entity.ImTopicChat, &entity.SubscribeMessage{
		Event: entity.SubEventContactStatus,
		Payload: jsonutil.Encode(entity.SubEventContactStatusPayload{
			Status: 2,
			UserId: client.Uid(),
		}),
	})
}
