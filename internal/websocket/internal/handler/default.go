package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/model"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/jwtutil"
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
	conn, err := im.NewConnect(ctx)
	if err != nil {
		logrus.Errorf("websocket connect error: %s", err.Error())
		return
	}

	// 创建客户端
	im.NewClient(ctx.Request.Context(), conn, &im.ClientOptions{
		Channel: im.Session.Default,
		Uid:     jwtutil.GetUid(ctx),
		Storage: c.cache,
	}, im.NewClientCallback(
		// 连接成功回调
		im.WithOpenCallback(func(client im.IClient) {
			c.open(client)
		}),
		// 接收消息回调
		im.WithMessageCallback(func(client im.IClient, message []byte) {
			c.message(client, message)
		}),
		// 关闭连接回调
		im.WithCloseCallback(func(client im.IClient, code int, text string) {
			c.close(client, code, text)
			// fmt.Printf("客户端[%d] 已关闭连接，关闭提示【%d】%s \n", client.ClientId(), code, text)
		}),
	))
}

// 连接成功回调事件
func (c *DefaultWebSocket) open(client im.IClient) {
	// 1.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.ClientUid())

	// 2.客户端加入群房间
	for _, id := range ids {
		_ = c.room.Add(context.Background(), &cache.RoomOption{
			Channel:  im.Session.Default.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      c.conf.ServerId(),
			Cid:      client.ClientId(),
		})
	}

	// 推送上线消息
	c.rds.Publish(context.Background(), entity.IMGatewayAll, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.ClientUid(),
			"status":  1,
		}),
	}))
}

// 消息接收回调事件
func (c *DefaultWebSocket) message(client im.IClient, message []byte) {
	content := string(message)

	event := gjson.Get(content, "event").String()

	switch event {

	// 对话键盘事件
	case entity.EventTalkKeyboard:
		var m *dto.KeyboardMessage
		if err := json.Unmarshal(message, &m); err == nil {
			c.rds.Publish(context.Background(), entity.IMGatewayAll, jsonutil.Encode(entity.MapStrAny{
				"event": entity.EventTalkKeyboard,
				"data": jsonutil.Encode(entity.MapStrAny{
					"sender_id":   m.Data.SenderID,
					"receiver_id": m.Data.ReceiverID,
				}),
			}))
		}

	// 对话消息读事件
	case entity.EventTalkRead:
		var m *dto.TalkReadMessage
		if err := json.Unmarshal(message, &m); err == nil {
			c.groupMemberService.Db().Model(&model.TalkRecords{}).Where("id in ? and receiver_id = ? and is_read = 0", m.Data.MsgIds, client.ClientUid()).Update("is_read", 1)

			c.rds.Publish(context.Background(), entity.IMGatewayAll, jsonutil.Encode(entity.MapStrAny{
				"event": entity.EventTalkRead,
				"data": jsonutil.Encode(entity.MapStrAny{
					"sender_id":   client.ClientUid(),
					"receiver_id": m.Data.ReceiverId,
					"ids":         m.Data.MsgIds,
				}),
			}))
		}
	default:
		fmt.Printf("消息事件未定义%s", event)
	}
}

// 客户端关闭回调事件
func (c *DefaultWebSocket) close(client im.IClient, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := c.groupMemberService.Dao().GetUserGroupIds(client.ClientUid())

	// 3.客户端退出群房间
	for _, id := range ids {
		_ = c.room.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Session.Default.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      c.conf.ServerId(),
			Cid:      client.ClientId(),
		})
	}

	// 推送下线消息
	c.rds.Publish(context.Background(), entity.IMGatewayAll, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.ClientUid(),
			"status":  0,
		}),
	}))
}
