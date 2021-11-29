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
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/service"
	"go-chat/config"
	"log"
	"strconv"
)

type DefaultWebSocket struct {
	rds                *redis.Client
	conf               *config.Config
	client             *service.ClientService
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
	handler := &DefaultWebSocket{rds: rds, conf: conf, client: client, room: room, groupMemberService: groupMemberService}

	im.Sessions.Default.SetCallbackHandler(handler)

	return handler
}

// Connect 初始化连接
func (ws *DefaultWebSocket) Connect(c *gin.Context) {
	conn, err := im.NewWebsocket(c)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}

	options := &im.ClientOption{
		Channel:       im.Sessions.Default,
		UserId:        auth.GetAuthUserID(c),
		ClientService: ws.client,
	}

	// 创建客户端
	im.NewClient(conn, options).InitConnection()
}

// Open 连接成功回调事件
func (ws *DefaultWebSocket) Open(client *im.Client) {
	// 1.查询用户群列表
	ids := ws.groupMemberService.GetUserGroupIds(client.Uid)

	// 2.客户端加入群房间
	for _, gid := range ids {
		_ = ws.room.Add(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name,
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(gid),
			Sid:      ws.conf.GetSid(),
			Cid:      client.ClientId,
		})
	}

	// 推送上线消息
	ws.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, jsonutil.JsonEncode(map[string]interface{}{
		"event_name": entity.EventOnlineStatus,
		"data": jsonutil.JsonEncode(map[string]interface{}{
			"user_id": client.Uid,
			"status":  1,
		}),
	}))
}

// Message 消息接收回调事件
func (ws *DefaultWebSocket) Message(message *im.ReceiveContent) {
	fmt.Printf("[%s]消息通知 Client:%d，Content: %s \n", message.Client.Channel.Name, message.Client.ClientId, message.Content)

	event := gjson.Get(message.Content, "event").String()

	switch event {
	case "event_keyboard":
		var m *wst.KeyboardMessage
		if err := json.Unmarshal([]byte(message.Content), &m); err == nil {
			ws.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, jsonutil.JsonEncode(map[string]interface{}{
				"event_name": entity.EventKeyboard,
				"data": jsonutil.JsonEncode(map[string]interface{}{
					"sender_id":   m.Data.SenderID,
					"receiver_id": m.Data.ReceiverID,
				}),
			}))
		}
	}
}

// Close 客户端关闭回调事件
func (ws *DefaultWebSocket) Close(client *im.Client, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := ws.groupMemberService.GetUserGroupIds(client.Uid)

	// 3.客户端退出群房间
	for _, gid := range ids {
		_ = ws.room.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Sessions.Default.Name,
			RoomType: entity.RoomGroupChat,
			Number:   strconv.Itoa(gid),
			Sid:      ws.conf.GetSid(),
			Cid:      client.ClientId,
		})
	}

	// 推送下线消息
	ws.rds.Publish(context.Background(), entity.SubscribeWsGatewayAll, jsonutil.JsonEncode(map[string]interface{}{
		"event_name": entity.EventOnlineStatus,
		"data": jsonutil.JsonEncode(map[string]interface{}{
			"user_id": client.Uid,
			"status":  0,
		}),
	}))
}
