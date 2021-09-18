package socket

// 响应的消息格式
type Message struct {
	Receiver []string `json:"-"`       // 接受者客户端 UUID
	IsAll    bool     `json:"-"`       // 是否推送所有客户端
	Event    string   `json:"event"`   // 消息事件
	Content  string   `json:"content"` // 推送信息
}
