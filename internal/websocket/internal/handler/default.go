package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/dto"
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
	return &DefaultWebSocket{rds: rds, conf: conf, cache: client, room: room, groupMemberService: groupMemberService}
}

// Connect 初始化连接
func (c *DefaultWebSocket) Connect(ctx *gin.Context) {
	conn, err := im.NewWebsocket(ctx)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", err.Error())
		return
	}

	// 创建客户端
	im.NewClient(conn, &im.ClientOptions{
		Channel: im.Sessions.Default,
		Uid:     jwt.GetUid(ctx),
		Storage: c.cache,

		// 设置回调方法
		CallBack: im.NewClientCallBack(im.WithClientCallBackOpen(func(client im.ClientInterface) {
			c.open(client)
		}), im.WithClientCallBackMessage(func(message *im.ReceiveContent) {
			c.message(message)
		}), im.WithClientCallBackClose(func(client im.ClientInterface, code int, text string) {
			c.close(client, code, text)
		})),
	})
}

// 连接成功回调事件
func (c *DefaultWebSocket) open(client im.ClientInterface) {
	// 1.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.ClientUid())

	// 2.客户端加入群房间
	for _, id := range ids {
		_ = c.room.Add(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(id),
			Sid:      c.conf.ServerId(),
			Cid:      client.ClientId(),
		})
	}

	// 推送上线消息
	c.rds.Publish(context.Background(), entity.IMGatewayAll, entity.JsonText{
		"event": entity.EventOnlineStatus,
		"data": entity.JsonText{
			"user_id": client.ClientUid(),
			"status":  1,
		}.Json(),
	}.Json())
}

// 消息接收回调事件
func (c *DefaultWebSocket) message(message *im.ReceiveContent) {
	// fmt.Printf("[%s]消息通知 Client:%d，Content: %s \n", message.Client.Channel().Name(), message.Client.ClientId(), message.Content)

	event := gjson.Get(message.Content, "event").String()

	switch event {
	case "event_keyboard":
		var m *dto.KeyboardMessage
		if err := json.Unmarshal([]byte(message.Content), &m); err == nil {
			c.rds.Publish(context.Background(), entity.IMGatewayAll, entity.JsonText{
				"event": entity.EventKeyboard,
				"data": entity.JsonText{
					"sender_id":   m.Data.SenderID,
					"receiver_id": m.Data.ReceiverID,
				}.Json(),
			}.Json())
		}
	}
}

// 客户端关闭回调事件
func (c *DefaultWebSocket) close(client im.ClientInterface, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.ClientUid())

	// 3.客户端退出群房间
	for _, id := range ids {
		_ = c.room.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name(),
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(id),
			Sid:      c.conf.ServerId(),
			Cid:      client.ClientId(),
		})
	}

	// 推送下线消息
	c.rds.Publish(context.Background(), entity.IMGatewayAll, entity.JsonText{
		"event": entity.EventOnlineStatus,
		"data": entity.JsonText{
			"user_id": client.ClientUid(),
			"status":  0,
		}.Json(),
	}.Json())
}
