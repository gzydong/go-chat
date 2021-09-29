package im

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"go-chat/app/service"
)

const (
	heartbeatCheckInterval int = 20 // 心跳检测时间
	heartbeatIdleTime      int = 50 // 心跳超时时间
)

// Client WebSocket 客户端连接信息
type Client struct {
	Conn          *websocket.Conn        // 客户端连接
	Uuid          string                 // 客户端唯一标识
	UserId        int                    // 用户ID
	LastTime      int64                  // 客户端最后心跳时间/心跳检测
	Channel       *ChannelManager        // 渠道分组
	ClientService *service.ClientService // 服务信息
}

// NewImClient ...
func NewImClient(conn *websocket.Conn, clientService *service.ClientService, userId int, channel *ChannelManager) *Client {
	client := &Client{
		Conn:          conn,
		Uuid:          uuid.NewV4().String(),
		UserId:        userId,
		LastTime:      time.Now().Unix(),
		Channel:       channel,
		ClientService: clientService,
	}

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("【%s】客户端关闭 %s | 关闭原因(%d): %s \n", client.Channel.Name, client.Uuid, code, text)

		channel.Handle.Close(client, code, text)

		client.Channel.RemoveClient(client)

		client.ClientService.UnBind(context.Background(), client.Channel.Name, client.Uuid)

		return nil
	})

	channel.RegisterClient(client)

	client.ClientService.Bind(context.Background(), channel.Name, client.Uuid, client.UserId)

	channel.Handle.Open(client)

	return client
}

// Close 关闭客户端连接
func (w *Client) Close(code int, message string) {
	// 触发客户端关闭回调事件
	Handler := w.Conn.CloseHandler()

	_ = Handler(code, message)

	if err := w.Conn.Close(); err != nil {
		log.Println("Close Error: ", err)
	}
}

// Heartbeat 心跳检测
func (w *Client) Heartbeat() {
	for {
		time.Sleep(time.Duration(heartbeatCheckInterval) * time.Second)

		if time.Now().Unix()-w.LastTime > int64(heartbeatIdleTime) {
			w.Close(2000, "心跳检测超时，连接自动关闭")
			break
		}
	}
}

// AcceptClient 接收客户端推送信息
func (w *Client) AcceptClient() {
	defer w.Conn.Close()

	for {
		// 读取ws中的数据
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "ping" {
			// 更新最后一次接受消息时间，用做心跳检测判断
			w.LastTime = time.Now().Unix()

			message = []byte("pong")

			// 写入ws数据
			err = w.Conn.WriteMessage(mt, message)
			if err != nil {
				break
			}

			continue
		}

		w.Channel.RecvChan <- &RecvMessage{
			Client:  w,
			Content: string(message),
		}
	}
}
