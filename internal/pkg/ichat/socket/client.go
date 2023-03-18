package socket

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"
	"go-chat/internal/pkg/jsonutil"
)

type IClient interface {
	Cid() int64                         // 客户端ID
	Uid() int                           // 客户端关联用户ID
	Close(code int, text string)        // 关闭客户端
	Write(data *ClientOutContent) error // 写入数据
	ChannelName() string
}

type IStorage interface {
	Bind(ctx context.Context, channel string, cid int64, uid int)
	UnBind(ctx context.Context, channel string, cid int64)
}

// ClientInContent 客户端接收消息体
type ClientInContent struct {
	IsAck   bool   // 是否需要 ack 回调
	Event   string // 消息事件
	Content []byte // 消息内容
}

// ClientOutContent 客户端输出的消息体
type ClientOutContent struct {
	AckId   string // ACK ID（唯一性）
	IsAck   bool   // 是否需要 ack 回调
	Retry   int    // 重试次数
	Content []byte // 消息内容
}

// Client WebSocket 客户端连接信息
type Client struct {
	conn     IConn                  // 客户端连接
	cid      int64                  // 客户端ID/客户端唯一标识
	uid      int                    // 用户ID
	lastTime int64                  // 客户端最后心跳时间/心跳检测
	channel  *Channel               // 渠道分组
	closed   int32                  // 客户端是否关闭连接
	storage  IStorage               // 缓存服务
	callBack ICallback              // 回调方法
	outChan  chan *ClientOutContent // 发送通道
}

type ClientOption struct {
	Uid     int      // 用户识别ID
	Channel *Channel // 渠道信息
	Storage IStorage // 自定义缓存组件，用于绑定用户与客户端的关系
	Buffer  int      // 缓冲区大小根据业务，自行调整
}

// NewClient 初始化客户端信息
func NewClient(ctx context.Context, conn IConn, opt *ClientOption, callBack ICallback) error {

	if opt.Buffer <= 0 {
		opt.Buffer = 10
	}

	if callBack == nil {
		panic("callBack is nil")
	}

	client := &Client{
		conn:     conn,
		cid:      Counter.GenID(),
		lastTime: time.Now().Unix(),
		uid:      opt.Uid,
		channel:  opt.Channel,
		storage:  opt.Storage,
		outChan:  make(chan *ClientOutContent, opt.Buffer),
		callBack: callBack,
	}

	// 设置客户端连接关闭回调事件
	conn.SetCloseHandler(client.close)

	// 绑定客户端映射关系
	if client.storage != nil {
		client.storage.Bind(ctx, client.channel.name, client.cid, client.uid)
	}

	// 注册客户端
	client.channel.addClient(client)

	// 触发自定义的 Open 事件
	client.callBack.Open(client)

	// 注册心跳管理
	health.addClient(client)

	return client.init()
}

// ChannelName Channel Name
func (c *Client) ChannelName() string {
	return c.channel.Name()
}

// Cid 获取客户端ID
func (c *Client) Cid() int64 {
	return c.cid
}

// Uid 获取客户端关联的用户ID
func (c *Client) Uid() int {
	return c.uid
}

// Close 关闭客户端连接
func (c *Client) Close(code int, message string) {
	defer c.conn.Close()

	// 触发客户端关闭回调事件
	if err := c.close(code, message); err != nil {
		log.Printf("[%s-%d-%d] client close err: %s \n", c.channel.Name(), c.cid, c.uid, err.Error())
	}
}

func (c *Client) Closed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

// Write 客户端写入数据
func (c *Client) Write(data *ClientOutContent) error {

	if c.Closed() {
		return fmt.Errorf("connection closed")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] [%s-%d-%d] client write err: %v \n", c.channel.Name(), c.cid, c.uid, err)
		}
	}()

	// 消息写入缓冲通道
	c.outChan <- data

	return nil
}

// 关闭回调
func (c *Client) close(code int, text string) error {

	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}

	close(c.outChan)

	// 触发连接关闭回调
	c.callBack.Close(c, code, text)

	if c.storage != nil {
		c.storage.UnBind(context.Background(), c.channel.name, c.cid)
	}

	// 渠道分组移除客户端
	c.channel.delClient(c)

	// 心跳管理移除客户端
	health.delClient(c)

	return nil
}

// 循环接收客户端推送信息
func (c *Client) loopAccept() {
	defer c.conn.Close()

	for {
		data, err := c.conn.Read()
		if err != nil {
			break
		}

		c.lastTime = time.Now().Unix()

		c.message(data)
	}
}

// 循环推送客户端信息
func (c *Client) loopWrite() {
	for data := range c.outChan {
		if c.Closed() {
			break
		}

		if err := c.conn.Write(data.Content); err != nil {
			log.Printf("[ERROR] [%s-%d-%d] client push write err: %v \n", c.channel.Name(), c.cid, c.uid, err)
			break
		}

		if data.IsAck && data.Retry > 0 {
			ack.add(data.AckId, &AckBufferBody{
				Cid:   c.cid,
				Uid:   int64(c.uid),
				Ch:    c.channel.name,
				Retry: data.Retry - 1,
				Body:  data.Content,
			})
		}
	}
}

func (c *Client) message(data []byte) {

	event := gjson.GetBytes(data, "event").String()
	if len(event) == 0 {
		return
	}

	switch event {
	// 心跳消息判断
	case "heartbeat", "event.heartbeat":
		_ = c.Write(&ClientOutContent{
			Content: jsonutil.Marshal(&Message{"heartbeat", "pong"}),
		})

	// 客户端 ACK 处理
	case "event.ack":
		// ack.remove("")
	default:
		// 触发消息回调
		c.callBack.Message(c, data)
	}
}

// 初始化连接
func (c *Client) init() error {

	// 推送心跳检测配置
	_ = c.Write(&ClientOutContent{
		Content: jsonutil.Marshal(&Message{
			Event: "connect",
			Content: map[string]any{
				"ping_interval": heartbeatInterval,
				"ping_timeout":  heartbeatTimeout,
			},
		}),
	})

	// 启动协程处理推送信息
	go c.loopWrite()

	go c.loopAccept()

	return nil
}
