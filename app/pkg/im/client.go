package im

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"go-chat/app/pkg/jsonutil"
	"time"

	"github.com/gorilla/websocket"
)

type StorageInterface interface {
	Bind(ctx context.Context, channel string, clientId string, id int)
	UnBind(ctx context.Context, channel string, clientId string)
}

// Client WebSocket 客户端连接信息
type Client struct {
	conn     *websocket.Conn  // 客户端连接
	cid      int64            // 客户端ID/客户端唯一标识
	uid      int              // 用户ID
	lastTime int64            // 客户端最后心跳时间/心跳检测
	channel  *Channel         // 渠道分组
	storage  StorageInterface // 缓存服务
	isClosed bool             // 客户端是否关闭连接
}

type ClientOptions struct {
	Uid     int
	Channel *Channel
	Storage StorageInterface
}

// NewClient 初始化客户端信息
func NewClient(conn *websocket.Conn, options *ClientOptions) *Client {
	client := &Client{
		conn:     conn,
		cid:      GenClientID.GetID(),
		lastTime: time.Now().Unix(),
		uid:      options.Uid,
		channel:  options.Channel,
		storage:  options.Storage,
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(func(code int, text string) error {
		client.isClosed = true

		options.Channel.handler.Close(client, code, text)

		options.Channel.delClient(client)

		client.storage.UnBind(context.Background(), client.Channel().name, fmt.Sprintf("%d", client.cid))

		// 通知心跳管理
		Heartbeat.delClient(client)

		return nil
	})

	// 绑定客户端映射关系
	client.storage.Bind(context.Background(), client.Channel().name, fmt.Sprintf("%d", client.cid), client.uid)

	// 注册客户端
	options.Channel.addClient(client)

	// 注册心跳管理
	Heartbeat.addClient(client)

	// 触发自定义的 open 事件
	options.Channel.handler.Open(client)

	return client
}

// ClientId 获取客户端ID
func (c *Client) ClientId() int64 {
	return c.cid
}

// Uid 获取客户端关联的用户ID
func (c *Client) Uid() int {
	return c.uid
}

// Channel 获取客户端通道信息
func (c *Client) Channel() *Channel {
	return c.channel
}

// IsClosed 判断客户端是否关闭连接
func (c *Client) IsClosed() bool {
	return c.isClosed
}

// Close 关闭客户端连接
func (c *Client) Close(code int, message string) {
	defer c.conn.Close()

	// 触发客户端关闭回调事件
	_ = c.conn.CloseHandler()(code, message)
}

// Write 客户端写入数据
func (c *Client) Write(messageType int, data []byte) error {

	if c.IsClosed() {
		return fmt.Errorf("client closed")
	}

	// 需要做线程安全处理
	return c.conn.WriteMessage(messageType, data)
}

// Init 初始化连接
func (c *Client) Init() {
	// 启动协程处理接收信息
	go c.accept()
}

// 循环接收客户端推送信息
func (c *Client) accept() {
	defer c.conn.Close()

	for {
		// 读取客户端中的数据
		mt, message, err := c.conn.ReadMessage()
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
			c.lastTime = time.Now().Unix()

			data, _ := jsonutil.JsonEncodeByte(&Message{"heartbeat", "pong"})

			_ = c.Write(mt, data)
			continue
		}

		if len(msg) > 0 {
			c.Channel().PushRecvChannel(&ReceiveContent{c, msg})
		}
	}
}
