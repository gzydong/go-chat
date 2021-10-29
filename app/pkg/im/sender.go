package im

type sender struct {
	broadcast bool     // 是否广播消息
	exclude   []int    // 排除的客户端
	receives  []int    // 推送的客户端
	message   *Message // 消息体
}

func NewSenderContent() *sender {
	return &sender{
		broadcast: false,
		exclude:   []int{},
		receives:  []int{},
	}
}

// SetBroadcast 设置广播推送
func (s *sender) SetBroadcast(value bool) *sender {
	s.broadcast = value
	return s
}

// SetMessage 设置推送数据
func (s *sender) SetMessage(msg *Message) *sender {
	s.message = msg
	return s
}

// SetReceive 设置推送客户端
func (s *sender) SetReceive(cid ...int) *sender {
	s.receives = append(s.receives, cid...)
	return s
}

// SetExclude 设置广播推送中需要过滤的客户端
func (s *sender) SetExclude(cid ...int) *sender {
	s.exclude = append(s.exclude, cid...)
	return s
}

// IsBroadcast 判断是否是广播推送
func (s *sender) IsBroadcast() bool {
	return s.broadcast
}

// SetMessage 设置推送数据
func (s *sender) GetMessage() interface{} {
	return s.message
}
