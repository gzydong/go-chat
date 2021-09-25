package im

// 响应的消息格式
type SendMessage struct {
	IsAll   bool     `json:"-"`       // 是否推送所有客户端
	Clients []string `json:"-"`       // 指定推送的客户端列表
	Event   string   `json:"event"`   // 消息事件
	Content string   `json:"content"` // 推送信息
}

// 接收的消息
type RecvMessage struct {
	Client  *Client // 接收的客户端
	Content string  // 接收的文本消息
}

// AddReceiver 添加接收者
func (m *SendMessage) AddReceiver(uuid string) {
	m.Clients = append(m.Clients, uuid)
}
