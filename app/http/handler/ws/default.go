package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"go-chat/app/cache"
	"go-chat/app/entity"
	wst "go-chat/app/http/dto/ws"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/im"
	"go-chat/app/service"
	"go-chat/config"
	"log"
	"strconv"
)

type DefaultWebSocket struct {
	rds                *redis.Client
	conf               *config.Config
	cache              *service.ClientService
	room               *cache.Room
	groupMemberService *service.GroupMemberService
}

func NewDefaultWebSocket(
	rds *redis.Client,
	conf *config.Config,
	client *service.ClientService,
	room *cache.Room,
	groupMemberService *service.GroupMemberService,
) *DefaultWebSocket {
	handler := &DefaultWebSocket{rds: rds, conf: conf, cache: client, room: room, groupMemberService: groupMemberService}

	im.Sessions.Default.SetHandler(handler)

	return handler
}

// Connect 初始化连接
func (c *DefaultWebSocket) Connect(ctx *gin.Context) {
	conn, err := im.NewWebsocket(ctx)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	// 创建客户端
	im.NewClient(conn, &im.ClientOptions{
		Channel: im.Sessions.Default,
		Uid:     auth.GetAuthUserID(ctx),
		Storage: c.cache,
	}).Init()
}

// Open 连接成功回调事件
func (c *DefaultWebSocket) Open(client *im.Client) {
	// 1.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.Uid())

	// 2.客户端加入群房间
	for _, gid := range ids {
		_ = c.room.Add(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(gid),
			Sid:      c.conf.GetSid(),
			Cid:      client.ClientId(),
		})
	}

	// 推送上线消息
	c.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, entity.JsonText{
		"event": entity.EventOnlineStatus,
		"data": entity.JsonText{
			"user_id": client.Uid(),
			"status":  1,
		}.Json(),
	}.Json())
}

// Message 消息接收回调事件
func (c *DefaultWebSocket) Message(message *im.ReceiveContent) {
	fmt.Printf("[%s]消息通知 Client:%d，Content: %s \n", message.Client.Channel().Name(), message.Client.ClientId(), message.Content)

	event := gjson.Get(message.Content, "event").String()

	switch event {
	case "event_keyboard":
		var m *wst.KeyboardMessage
		if err := json.Unmarshal([]byte(message.Content), &m); err == nil {
			c.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, entity.JsonText{
				"event": entity.EventKeyboard,
				"data": entity.JsonText{
					"sender_id":   m.Data.SenderID,
					"receiver_id": m.Data.ReceiverID,
				}.Json(),
			}.Json())
		}
	}
}

// Close 客户端关闭回调事件
func (c *DefaultWebSocket) Close(client *im.Client, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.Uid())

	// 3.客户端退出群房间
	for _, gid := range ids {
		_ = c.room.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(gid),
			Sid:      c.conf.GetSid(),
			Cid:      client.ClientId(),
		})
	}

	// 推送下线消息
	c.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, entity.JsonText{
		"event": entity.EventOnlineStatus,
		"data": entity.JsonText{
			"user_id": client.Uid(),
			"status":  0,
		}.Json(),
	}.Json())
}
