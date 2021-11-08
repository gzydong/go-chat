package im

// Message 客户端交互的消息体
type Message struct {
	Event   string      `json:"event"`   // 事件名称
	Content interface{} `json:"content"` // 消息内容
}

// ReceiveContent 接收的消息
type ReceiveContent struct {
	Client  *Client // 接收的客户端
	Content string  // 接收的文本消息
}

// SenderContent 推送的消息
type SenderContent struct {
	broadcast bool     // 是否广播消息
	exclude   []int    // 排除的用户
	receives  []int    // 推送的用户
	message   *Message // 消息体
}

func NewSenderContent() *SenderContent {
	return &SenderContent{
		broadcast: false,
		exclude:   make([]int, 0),
		receives:  make([]int, 0),
	}
}

// SetBroadcast 设置广播推送
func (s *SenderContent) SetBroadcast(value bool) *SenderContent {
	s.broadcast = value
	return s
}

// SetMessage 设置推送数据
func (s *SenderContent) SetMessage(msg *Message) *SenderContent {
	s.message = msg
	return s
}

// SetReceive 设置推送客户端
func (s *SenderContent) SetReceive(uid ...int) *SenderContent {
	s.receives = append(s.receives, uid...)
	return s
}

// SetExclude 设置广播推送中需要过滤的客户端
func (s *SenderContent) SetExclude(uid ...int) *SenderContent {
	s.exclude = append(s.exclude, uid...)
	return s
}

// IsBroadcast 判断是否是广播推送
func (s *SenderContent) IsBroadcast() bool {
	return s.broadcast
}

// GetMessage 获取消息内容
func (s *SenderContent) GetMessage() interface{} {
	return s.message
}
