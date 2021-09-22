package im

// 响应的消息格式
type Message struct {
	IsAll   bool     `json:"-"`       // 是否推送所有客户端
	Clients []string `json:"-"`       // 指定推送的客户端列表
	Event   string   `json:"event"`   // 消息事件
	Content string   `json:"content"` // 推送信息
}

// AddReceiver 添加接收者
func (m *Message) AddReceiver(uuid string) {
	m.Clients = append(m.Clients, uuid)
}
