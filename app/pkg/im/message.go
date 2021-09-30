package im

// WebSocketMessage WebSocket 消息格式
type WebSocketMessage struct {
	Type    string      `json:"type"`    // 消息类型
	Content interface{} `json:"content"` // 推送信息
}

// SendMessage 响应的消息格式
type SendMessage struct {
	IsAll   bool   `json:"-"`       // 是否推送所有客户端
	Clients []int  `json:"-"`       // 指定推送的客户端列表
	Event   string `json:"event"`   // 消息事件
	Content string `json:"content"` // 推送信息
}

// RecvMessage 接收的消息
type RecvMessage struct {
	Client  *Client // 接收的客户端
	Content string  // 接收的文本消息
}

// AddReceiver 添加接收者
func (m *SendMessage) AddReceiver(clientId int) {
	m.Clients = append(m.Clients, clientId)
}
