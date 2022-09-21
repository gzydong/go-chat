package event

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/dto"
)

type DefaultEvent struct {
	redis         *redis.Client
	config        *config.Config
	cache         *service.ClientService
	roomStorage   *cache.RoomStorage
	memberService *service.GroupMemberService
}

func NewDefaultEvent(rds *redis.Client, conf *config.Config, cache *service.ClientService, room *cache.RoomStorage, groupMemberService *service.GroupMemberService) *DefaultEvent {
	return &DefaultEvent{redis: rds, config: conf, cache: cache, roomStorage: room, memberService: groupMemberService}
}

func (d *DefaultEvent) OnOpen(client im.IClient) {
	// 1.查询用户群列表
	ids := d.memberService.Dao().GetUserGroupIds(client.Uid())

	// 2.客户端加入群房间
	for _, id := range ids {
		_ = d.roomStorage.Add(context.Background(), &cache.RoomOption{
			Channel:  im.Session.Default.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      d.config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	// 推送上线消息
	d.redis.Publish(context.Background(), entity.ImTopicDefault, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.Uid(),
			"status":  1,
		}),
	}))
}

func (d *DefaultEvent) OnMessage(client im.IClient, message []byte) {
	content := string(message)

	event := gjson.Get(content, "event").String()

	switch event {

	// 对话键盘事件
	case entity.EventTalkKeyboard:
		var m *dto.KeyboardMessage
		if err := json.Unmarshal(message, &m); err == nil {
			d.redis.Publish(context.Background(), entity.ImTopicDefault, jsonutil.Encode(entity.MapStrAny{
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
			d.memberService.Db().Model(&model.TalkRecords{}).Where("id in ? and receiver_id = ? and is_read = 0", m.Data.MsgIds, client.Uid()).Update("is_read", 1)

			d.redis.Publish(context.Background(), entity.ImTopicDefault, jsonutil.Encode(entity.MapStrAny{
				"event": entity.EventTalkRead,
				"data": jsonutil.Encode(entity.MapStrAny{
					"sender_id":   client.Uid(),
					"receiver_id": m.Data.ReceiverId,
					"ids":         m.Data.MsgIds,
				}),
			}))
		}
	default:
		fmt.Printf("消息事件未定义%s", event)
	}
}

func (d *DefaultEvent) OnClose(client im.IClient, code int, text string) {
	// 1.判断用户是否是多点登录

	// 2.查询用户群列表
	ids := d.memberService.Dao().GetUserGroupIds(client.Uid())

	// 3.客户端退出群房间
	for _, id := range ids {
		_ = d.roomStorage.Del(context.Background(), &cache.RoomOption{
			Channel:  im.Session.Default.Name(),
			RoomType: entity.RoomImGroup,
			Number:   strconv.Itoa(id),
			Sid:      d.config.ServerId(),
			Cid:      client.Cid(),
		})
	}

	// 推送下线消息
	d.redis.Publish(context.Background(), entity.ImTopicDefault, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventOnlineStatus,
		"data": jsonutil.Encode(entity.MapStrAny{
			"user_id": client.Uid(),
			"status":  0,
		}),
	}))
}
