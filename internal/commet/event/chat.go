package event

import (
	"context"
	"log"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"go-chat/config"
	"go-chat/internal/commet/event/chat"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type ChatEvent struct {
	Redis           *redis.Client
	Config          *config.Config
	RoomStorage     *cache.RoomStorage
	GroupMemberRepo *repo.GroupMember

	MemberService service.IGroupMemberService
	Handler       *chat.Handler
}

// OnOpen 连接成功回调事件
func (c *ChatEvent) OnOpen(client socket.IClient) {

	ctx := context.TODO()

	// 1.查询用户群列表
	ids := c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid())

	// 2.客户端加入群房间
	rooms := make([]*cache.RoomOption, 0, len(ids))
	for _, id := range ids {
		rooms = append(rooms, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      c.Config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	if err := c.RoomStorage.BatchAdd(ctx, rooms); err != nil {
		log.Println("加入群聊失败", err.Error())
	}

	// 推送上线消息
	c.Redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventContactStatus,
		"data": jsonutil.Encode(map[string]any{
			"user_id": client.Uid(),
			"status":  1,
		}),
	}))
}

// OnMessage 消息回调事件
func (c *ChatEvent) OnMessage(client socket.IClient, message []byte) {
	val, err := sonic.Get(message, "event")
	if err != nil {
		return
	}

	// 获取事件名
	event, _ := val.String()
	if event != "" {
		// 触发事件
		c.Handler.Call(context.TODO(), client, event, message)
	}
}

// OnClose 连接关闭回调事件
func (c *ChatEvent) OnClose(client socket.IClient, code int, text string) {

	log.Println("client close: ", client.Uid(), client.Cid(), client.Channel().Name(), code, text)

	ctx := context.TODO()

	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := c.GroupMemberRepo.GetUserGroupIds(ctx, client.Uid())

	// 3.客户端退出群房间
	rooms := make([]*cache.RoomOption, 0, len(ids))
	for _, id := range ids {
		rooms = append(rooms, &cache.RoomOption{
			Channel:  socket.Session.Chat.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      c.Config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	if err := c.RoomStorage.BatchDel(ctx, rooms); err != nil {
		log.Println("退出群聊失败", err.Error())
	}

	// 推送下线消息
	c.Redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventContactStatus,
		"data": jsonutil.Encode(map[string]any{
			"user_id": client.Uid(),
			"status":  0,
		}),
	}))
}
