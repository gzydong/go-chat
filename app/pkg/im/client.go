package im

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"go-chat/app/pkg/jsonutil"
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
	ClientId      int64                  // 客户端ID/客户端唯一标识
	Uid           int                    // 用户ID
	LastTime      int64                  // 客户端最后心跳时间/心跳检测
	Channel       *ChannelManage         // 渠道分组
	ClientService *service.ClientService // 服务信息
	IsClosed      bool                   // 客户端是否关闭连接
}

type ClientOption struct {
	UserId        int
	Channel       *ChannelManage
	ClientService *service.ClientService
}

// NewClient ...
func NewClient(conn *websocket.Conn, options *ClientOption) *Client {
	client := &Client{
		Conn:          conn,
		ClientId:      GenClientID.GetID(),
		LastTime:      time.Now().Unix(),
		Uid:           options.UserId,
		Channel:       options.Channel,
		ClientService: options.ClientService,
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(func(code int, text string) error {
		client.IsClosed = true

		client.Channel.Handler.Close(client, code, text)

		client.Channel.RemoveClient(client)

		client.ClientService.UnBind(context.Background(), client.Channel.Name, fmt.Sprintf("%d", client.ClientId))

		return nil
	})

	// 注册客户端
	client.Channel.RegisterClient(client)

	// 绑定客户端映射关系
	client.ClientService.Bind(context.Background(), client.Channel.Name, fmt.Sprintf("%d", client.ClientId), client.Uid)

	// 触发自定义的 open 事件
	client.Channel.Handler.Open(client)

	return client
}

// Close 关闭客户端连接
func (w *Client) Close(code int, message string) {
	defer w.Conn.Close()

	// 触发客户端关闭回调事件
	handler := w.Conn.CloseHandler()

	_ = handler(code, message)
}

// Write 客户端写入数据
func (w *Client) Write(messageType int, data []byte) error {
	// 需要做线程安全处理
	return w.Conn.WriteMessage(messageType, data)
}

// InitConnection 初始化连接
func (w *Client) InitConnection() {
	// 启动协程处理接收信息
	go w.accept()

	// 启动客户端心跳检测
	go w.heartbeat()
}

// 心跳检测
func (w *Client) heartbeat() {
	for {
		<-time.After(heartbeatCheckInterval)

		if w.IsClosed {
			break
		}

		if int(time.Now().Unix()-w.LastTime) > heartbeatIdleTime {
			w.Close(2000, "心跳检测超时，连接自动关闭")
			break
		}
	}
}

// 循环接收客户端推送信息
func (w *Client) accept() {
	defer w.Conn.Close()

	for {
		// 读取客户端中的数据
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			break
		}

		msg := string(message)

		res := gjson.Get(msg, "event")

		// 判断消息格式是否正确
		if !res.Exists() {
			continue
		}

		// 心跳消息判断
		if res.String() == "heartbeat" {
			w.LastTime = time.Now().Unix()
			_ = w.Write(mt, []byte(jsonutil.JsonEncode(&Message{
				Event:   "heartbeat",
				Content: "pong",
			})))
			continue
		}

		if len(msg) > 0 {
			w.Channel.PushRecvChannel(&ReceiveContent{
				Client:  w,
				Content: msg,
			})
		}
	}
}
