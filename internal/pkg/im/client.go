package im

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"go-chat/internal/pkg/jsonutil"
)

type ClientInterface interface {
	ClientId() int64                // 获取客户端ID
	ClientUid() int                 // 获取客户端关联用户ID
	IsClosed() bool                 // 判断客户端是否关闭
	Close(code int, message string) // 关闭客户端
	Write(data []byte) error        // 客户端写入数据
}

type StorageInterface interface {
	Bind(ctx context.Context, channel string, clientId string, id int)
	UnBind(ctx context.Context, channel string, clientId string)
}

// ClientInContent 客户端接收消息体
type ClientInContent struct {
	IsAck   bool   // 是否需要 ack 回调
	Event   string // 消息事件
	Content []byte // 消息内容
}

// ClientOutContent 客户端输出的消息体
type ClientOutContent struct {
	IsAck   bool   // 是否需要 ack 回调
	Retry   int    // 重试次数
	Content []byte // 消息内容
}

// Client WebSocket 客户端连接信息
type Client struct {
	conn     *websocket.Conn         // 客户端连接
	cid      int64                   // 客户端ID/客户端唯一标识
	uid      int                     // 用户ID
	lastTime int64                   // 客户端最后心跳时间/心跳检测
	channel  *Channel                // 渠道分组
	storage  StorageInterface        // 缓存服务
	isClosed bool                    // 客户端是否关闭连接
	outChan  chan []byte             // 发送通道
	callBack ClientCallBackInterface // 回调方法
}

type ClientOptions struct {
	Uid      int
	Channel  *Channel
	Storage  StorageInterface
	CallBack ClientCallBackInterface // 回调方法设置
}

// NewClient 初始化客户端信息
func NewClient(conn *websocket.Conn, opt *ClientOptions, callBack ClientCallBackInterface) ClientInterface {
	client := &Client{
		conn:     conn,
		cid:      Counter.GetID(),
		lastTime: time.Now().Unix(),
		uid:      opt.Uid,
		channel:  opt.Channel,
		storage:  opt.Storage,
		outChan:  make(chan []byte, 5), // 缓冲区大小根据业务，自行调整
		callBack: callBack,
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(client.setCloseHandler)

	if client.storage != nil {
		// 绑定客户端映射关系
		client.storage.Bind(context.Background(), opt.Channel.name, fmt.Sprintf("%d", client.cid), client.uid)
	}

	// 注册客户端
	client.channel.addClient(client)

	// 触发自定义的 open 事件
	client.callBack.Open(client)

	// 注册心跳管理
	Heartbeat.addClient(client)

	// 推送心跳检测配置
	_ = client.Write(jsonutil.EncodeByte(&Message{
		Event: "connect",
		Content: map[string]interface{}{
			"ping_interval": HeartbeatInterval,
			"ping_timeout":  HeartbeatTimeout,
		},
	}))

	return client.init()
}

// ClientId 获取客户端ID
func (c *Client) ClientId() int64 {
	return c.cid
}

// ClientUid 获取客户端关联的用户ID
func (c *Client) ClientUid() int {
	return c.uid
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
func (c *Client) Write(data []byte) error {

	if c.IsClosed() {
		return fmt.Errorf("websocket client closed")
	}

	// 消息写入缓冲通道
	c.outChan <- data

	return nil
}

// 关闭回调
func (c *Client) setCloseHandler(code int, text string) error {
	if !c.isClosed {
		close(c.outChan) // 关闭通道
	}

	c.isClosed = true

	// 触发连接关闭回调
	c.callBack.Close(c, code, text)

	// 解绑关联
	if c.storage != nil {
		c.storage.UnBind(context.Background(), c.channel.name, fmt.Sprintf("%d", c.cid))
	}

	// 渠道分组移除客户端
	c.channel.delClient(c)

	// 心跳管理移除客户端
	Heartbeat.delClient(c)

	return nil
}

// 循环接收客户端推送信息
func (c *Client) loopAccept() {
	defer c.conn.Close()

	for {
		// 读取客户端中的数据
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		msg := string(message)

		res := gjson.Get(msg, "event")

		// 判断消息格式是否正确
		if !res.Exists() {
			continue
		}

		switch res.String() {
		case "heartbeat": // 心跳消息判断
			c.lastTime = time.Now().Unix()

			_ = c.Write(jsonutil.EncodeByte(&Message{"heartbeat", "pong"}))
		case "ack":
			ackManage.Del(&AckBufferOption{
				Client: c,
				MsgID:  "",
			})
		default:
			// 触发消息回调
			c.callBack.Message(&ReceiveContent{c, msg})
		}
	}
}

// 循环推送客户端信息
func (c *Client) loopWrite() {
	for msg := range c.outChan {

		if c.isClosed {
			break
		}

		_ = c.conn.WriteMessage(websocket.TextMessage, msg)

		// 这里需要消息推送 ack 通道
	}
}

// Init 初始化连接
func (c *Client) init() *Client {
	// 启动协程处理接收信息
	go c.loopAccept()

	// 启动协程处理推送信息
	go c.loopWrite()

	return c
}
