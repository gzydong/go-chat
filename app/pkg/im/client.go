package im

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"go-chat/app/service"
)

const (
	heartbeatCheckInterval = 20 * time.Second // 心跳检测时间
	heartbeatIdleTime      = 50               // 心跳超时时间
)

// Client WebSocket 客户端连接信息
type Client struct {
	Conn          *websocket.Conn        // 客户端连接
	ClientId      int                    // 客户端ID/客户端唯一标识
	UserId        int                    // 用户ID
	LastTime      int64                  // 客户端最后心跳时间/心跳检测
	Channel       *ChannelManager        // 渠道分组
	ClientService *service.ClientService // 服务信息
}

// NewClient ...
func NewClient(conn *websocket.Conn, clientService *service.ClientService, userId int, channel *ChannelManager) *Client {
	client := &Client{
		Conn:          conn,
		ClientId:      NewClientID(),
		UserId:        userId,
		LastTime:      time.Now().Unix(),
		Channel:       channel,
		ClientService: clientService,
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(func(code int, text string) error {
		channel.Handle.Close(client, code, text)

		client.Channel.RemoveClient(client)

		client.ClientService.UnBind(context.Background(), client.Channel.Name, strconv.Itoa(client.ClientId))

		return nil
	})

	// 注册客户端
	channel.RegisterClient(client)

	// 绑定客户端映射关系
	client.ClientService.Bind(context.Background(), channel.Name, strconv.Itoa(client.ClientId), client.UserId)

	// 触发自定义的 open 事件
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

// heartbeat 心跳检测
func (w *Client) heartbeat() {
	for {
		<-time.After(heartbeatCheckInterval)

		if int(time.Now().Unix()-w.LastTime) > heartbeatIdleTime {
			w.Close(2000, "心跳检测超时，连接自动关闭")
			break
		}
	}
}

// accept 接收客户端推送信息
func (w *Client) accept() {
	defer w.Close(3000, "[协程异常] AcceptClient 已结束")

	for {
		// 读取客户端中的数据
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			break
		}

		msg := string(message)

		// 心跳消息判断
		if msg == "ping" {
			w.LastTime = time.Now().Unix()

			if w.Conn.WriteMessage(websocket.PongMessage, []byte("pong")) != nil {
				break
			}

			continue
		}

		// todo 这里需要验证消息格式，未知格式直接忽略

		if len(msg) > 0 {
			w.Channel.RecvMessage(&RecvMessage{
				Client:  w,
				Content: msg,
			})
		}
	}
}

func (w *Client) Start() {
	// 启动协程处理接收信息
	go w.accept()

	// 启动客户端心跳检测
	go w.heartbeat()
}
