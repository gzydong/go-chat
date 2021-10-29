package im

import (
	"reflect"
)

type sender struct {
	channel   string   // 渠道名
	broadcast bool     // 是否广播消息
	exclude   []int    // 排除的客户端
	receive   []int    // 推送的客户端
	message   *Message // 消息体
}

type Message struct {
	Event   string      // 事件名称
	Content interface{} // 消息内容
}

func NewSender(channel string) *sender {
	return &sender{
		channel:   channel,
		broadcast: false,
		exclude:   []int{},
		receive:   []int{},
	}
}

// SetBroadcast 设置广播推送
func (s *sender) SetBroadcast(value bool) {
	s.broadcast = value
}

// SetMessage 设置推送数据
func (s *sender) SetMessage(msg *Message) {
	s.message = msg
}

func (s *sender) SetReceive(cid ...int) {
	s.receive = append(s.receive, cid...)
}

func (s *sender) SetExclude(cid ...int) {
	s.exclude = append(s.exclude, cid...)
}

func (s *sender) GetReceive() []int {
	return s.receive
}

func (s *sender) GetExclude() []int {
	return s.exclude
}

// Send 发送数据
func (s *sender) Send() {
	num := reflect.TypeOf(Manager).Elem().NumField()
	for i := 0; i < num; i++ {
		group := reflect.ValueOf(Manager).Elem().Field(i).Interface()

		if reflect.ValueOf(group).Elem().FieldByName("Name").String() == s.channel {
			// reflect.ValueOf(group).Elem()
		}
	}
}
