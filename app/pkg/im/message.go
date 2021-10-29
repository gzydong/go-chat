package im

// Message 客户端交互的消息体
type Message struct {
	Event   string      `json:"event"`   // 事件名称
	Content interface{} `json:"content"` // 消息内容
}

// RecvMessage 接收的消息
type RecvMessage struct {
	Client  *Client // 接收的客户端
	Content string  // 接收的文本消息
}
