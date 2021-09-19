package im

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	heartbeatCheckInterval int = 20 // 心跳检测时间
	heartbeatIdleTime      int = 50 // 心跳超时时间
)

// Client WebSocket 客户端连接信息
type Client struct {
	Conn     *websocket.Conn // 客户端连接
	Uuid     string          // 客户端唯一标识
	UserId   int             // 用户ID
	LastTime int64           // 客户端最后心跳时间/心跳检测
	Channel  string          // 渠道分组
}

type CloseFunc func(c *Client) bool

// NewImClient ...
func NewImClient(conn *websocket.Conn, userId int, Channel *ChannelManager) *Client {
	client := &Client{
		Conn:     conn,
		Uuid:     uuid.NewV4().String(),
		UserId:   userId,
		LastTime: time.Now().Unix(),
		Channel:  Channel.Name,
	}

	Channel.RegisterClient(client)

	return client
}

// Close 关闭客户端连接
func (w *Client) Close(code int, message string) {
	// 触发客户端关闭回调事件
	Handler := w.Conn.CloseHandler()

	_ = Handler(code, message)

	w.Conn.Close()
}

// Heartbeat 心跳检测
func (w *Client) Heartbeat(fn CloseFunc) {
	for {
		time.Sleep(time.Duration(heartbeatCheckInterval) * time.Second)

		if time.Now().Unix()-w.LastTime > int64(heartbeatIdleTime) {
			isOk := fn(w)
			if isOk {
				w.Close(500, "心跳检测超时，连接自动关闭")
				break
			}
		}
	}
}

func (w *Client) AcceptClient() {
	defer w.Conn.Close()

	for {
		//读取ws中的数据
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 更新最后一次接受消息时间，用做心跳检测判断
		w.LastTime = time.Now().Unix()

		if string(message) == "ping" {
			message = []byte("pong")

			//写入ws数据
			err = w.Conn.WriteMessage(mt, message)
			if err != nil {
				break
			}

			continue
		}
	}
}

// SetCloseHandler 设置客户端关闭回调处理事件
func (w *Client) SetCloseHandler(fn func(code int, text string) error) {
	w.Conn.SetCloseHandler(func(code int, text string) error {
		_ = fn(code, text)

		//el := reflect.ValueOf(Manager).Elem()
		//for i := 0; i < el.NumField(); i++ {
		//	if w.Channel == el.Field(i).Elem().FieldByName("Name").String() {
		//
		//		params := make([]reflect.Value, 1)
		//		params[0] = reflect.ValueOf(w)
		//
		//		el.Field(i).Elem().MethodByName("RemoveClient").Call(params)
		//
		//		break
		//	}
		//}

		return nil
	})
}
