package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/pkg/core/consumer"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/server"
)

var _ consumer.IConsumerHandle = (*RoomControl)(nil)

// RoomControl 群成员进群退群通知
type RoomControl struct {
	Room *socket.RoomStorage
}

func (g *RoomControl) Topic() string {
	return "im.room.control"
}

func (g *RoomControl) Channel() string {
	return server.ID()
}

func (g *RoomControl) Touch() bool {
	return false
}

type GroupControlMessage struct {
	Action    int32   `json:"action"`    // 操作方式 1:入群 2:退群 3:解散群
	GroupId   int32   `json:"group_id"`  // 群ID
	UserIds   []int32 `json:"user_ids"`  // 成员ID
	Timestamp int64   `json:"timestamp"` // 操作时间
}

func (g *RoomControl) Do(ctx context.Context, message []byte, attempts uint16) error {
	var data GroupControlMessage
	if err := json.Unmarshal(message, &data); err != nil {
		fmt.Printf("RoomControl Unmarshal err:%v", err)
		return err
	}

	switch data.Action {
	case 1, 2:
		if len(data.UserIds) == 0 {
			fmt.Println("RoomControl UserIds is empty")
			return nil
		}

		if data.Action == 1 {
			err := g.Room.BatchInsert(data.GroupId, []int64{}, data.Timestamp)
			if err != nil {
				fmt.Println("RoomControl BatchInsert err:", err)
			}
		} else {
			err := g.Room.BatchDelete(data.GroupId, []int64{}, data.Timestamp)
			if err != nil {
				fmt.Println("RoomControl BatchDelete err:", err)
			}
		}
	case 3:
		err := g.Room.DeleteRoom(data.GroupId)
		if err != nil {
			fmt.Println("RoomControl DeleteRoom err:", err)
		}
	}

	return nil
}
