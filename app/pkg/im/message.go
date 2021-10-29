package im

// Message 客户端交互的消息体
type Message struct {
	Event   string      `json:"event"`   // 事件名称
	Content interface{} `json:"content"` // 消息内容
}

// ClientContent 接收的消息
type ClientContent struct {
	Client  *Client // 接收的客户端
	Content string  // 接收的文本消息
}

type SenderContent struct {
	broadcast bool     // 是否广播消息
	exclude   []int    // 排除的客户端
	receives  []int    // 推送的客户端
	message   *Message // 消息体
}

func NewSenderContent() *SenderContent {
	return &SenderContent{
		broadcast: false,
		exclude:   []int{},
		receives:  []int{},
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
func (s *SenderContent) SetReceive(cid ...int) *SenderContent {
	s.receives = append(s.receives, cid...)
	return s
}

// SetExclude 设置广播推送中需要过滤的客户端
func (s *SenderContent) SetExclude(cid ...int) *SenderContent {
	s.exclude = append(s.exclude, cid...)
	return s
}

// IsBroadcast 判断是否是广播推送
func (s *SenderContent) IsBroadcast() bool {
	return s.broadcast
}

// SetMessage 设置推送数据
func (s *SenderContent) GetMessage() interface{} {
	return s.message
}
